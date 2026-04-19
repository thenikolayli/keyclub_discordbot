package eventutils

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/internal"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/sheets/v4"
)

func SyncEvents(ctx context.Context, app *internal.App) error {
	if !app.EventSync.ShouldSync() {
		remaining := app.EventSync.UpdateTimeout - time.Since(app.EventSync.LastUpdated)
		return fmt.Errorf("Not enough time has passed since the last update, wait %v more seconds.", remaining.Seconds())
	}
	app.EventSync.Mutex.Lock()
	defer app.EventSync.Mutex.Unlock()

	if err := syncEventsFromCalendar(ctx, app); err != nil {
		return err
	}
	if err := syncEventsFromSheet(ctx, app); err != nil {
		return err
	}
	app.EventSync.LastUpdated = time.Now()
	return nil
}

func syncEventsFromCalendar(ctx context.Context, app *internal.App) error {
	calendarEvents, err := app.GoogleServices.Calendar.Events.List(app.Config.CalendarID).TimeMin(time.Now().Format(time.RFC3339)).TimeMax(time.Now().AddDate(0, 1, 0).Format(time.RFC3339)).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Failed to get events from calendar: %w", err)
	}

	transaction, err := app.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to create a transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			transaction.Rollback()
		}
	}()

	var wg sync.WaitGroup
	var upsertMutex sync.Mutex
	var workErr error

	for _, calEvent := range calendarEvents.Items {
		if len(calEvent.Attachments) == 0 {
			continue
		}
		wg.Add(1)
		go func(ev *calendar.Event) {
			defer wg.Done()
			formattedEvent, _, err := GetEventInfo(ctx, DocsUrlToId(ev.Attachments[0].FileUrl), app.GoogleServices.Docs)
			if err != nil {
				fmt.Printf("Failed to get event info for event %s: %v\n", ev.Summary, err)
				return
			}
			upsertMutex.Lock()
			defer upsertMutex.Unlock()
			if workErr != nil {
				return
			}
			if err := upsertEvent(ctx, formattedEvent, transaction); err != nil {
				workErr = err
			}
		}(calEvent)
	}
	wg.Wait()
	if workErr != nil {
		return workErr
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	committed = true
	return nil
}

// syncs the events database entries from the events spreadsheet
// fetches values via an api call to the events spreadsheet
// formats the response to event structs
// updates the database based on structs
func syncEventsFromSheet(ctx context.Context, app *internal.App) error {
	eventValueRanges, err := getEventValueRanges(ctx, app)
	if err != nil {
		return fmt.Errorf("Failed to update events: %w", err)
	}
	formattedEventStructs := getEventStructs(eventValueRanges)

	transaction, err := app.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to create a transaction: %w", err)
	}
	for _, each := range formattedEventStructs {
		err := upsertEvent(ctx, each, transaction)
		if err != nil {
			return err
		}
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	app.EventSync.LastUpdated = time.Now()
	return nil
}

// fetches and returns google sheets api value ranges (unformatted)
func getEventValueRanges(ctx context.Context, app *internal.App) ([]*sheets.ValueRange, error) {
	r := app.Config.EventsSheetRanges
	data, err := app.GoogleServices.Sheets.Spreadsheets.Values.BatchGet(app.Config.SpreadsheetID).Ranges(
		r.Events,
		r.Dates,
		r.StartTimes,
		r.EndTimes,
		r.Addresses,
		r.NofSlots,
		r.NofVolunteers,
		r.TotalHours,
		r.Leaders,
		r.MadeBy,
	).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to batch get spreadsheet ranges: %v", err)
	}
	return data.ValueRanges, nil
}

// takes the api call value ranges and turns them into an array of event structs
func getEventStructs(eventValueRanges []*sheets.ValueRange) []Event {
	// gets length based on the length of the events column
	eventValueRangesLength := len(eventValueRanges[0].Values)
	formattedEventArray := make([]Event, eventValueRangesLength)

	normalizedEvents := genericutils.NormalizeStringValues(eventValueRanges[0].Values, eventValueRangesLength)
	normalizedDates := genericutils.NormalizeStringValues(eventValueRanges[1].Values, eventValueRangesLength)
	normalizedStartTimes := genericutils.NormalizeStringValues(eventValueRanges[2].Values, eventValueRangesLength)
	normalizedEndTimes := genericutils.NormalizeStringValues(eventValueRanges[3].Values, eventValueRangesLength)
	normalizedAddresses := genericutils.NormalizeStringValues(eventValueRanges[4].Values, eventValueRangesLength)
	normalizedNofSlots := genericutils.NormalizeIntValues(eventValueRanges[5].Values, eventValueRangesLength)
	normalizedNofVolunteers := genericutils.NormalizeIntValues(eventValueRanges[6].Values, eventValueRangesLength)
	normalizedTotalHours := genericutils.NormalizeFloatValues(eventValueRanges[7].Values, eventValueRangesLength)
	normalizedLeaders := genericutils.NormalizeStringValues(eventValueRanges[8].Values, eventValueRangesLength)
	normalizedMadeBy := genericutils.NormalizeStringValues(eventValueRanges[9].Values, eventValueRangesLength)

	for i := range eventValueRangesLength {
		formattedEventArray[i] = Event{
			Name:          normalizedEvents[i],
			Date:          normalizedDates[i],
			StartTime:     normalizedStartTimes[i],
			EndTime:       normalizedEndTimes[i],
			Address:       normalizedAddresses[i],
			NofSlots:      normalizedNofSlots[i],
			NofVolunteers: normalizedNofVolunteers[i],
			TotalHours:    normalizedTotalHours[i],
			Leaders:       normalizedLeaders[i],
			MadeBy:        normalizedMadeBy[i],
			SignUpUrl:     "", // sheet doesn't have this column
			ID:            -1,
		}
	}

	return formattedEventArray
}

// takes in a formatted event struct and a transaction and upserts their row
// checks if an event with the same name exists
// if it doesn't, insert it
// otherwise, update
func upsertEvent(ctx context.Context, event Event, transaction *sqlx.Tx) error {
	result := Event{}
	err := transaction.GetContext(
		ctx, &result,
		"SELECT * from events WHERE name = ? LIMIT 1",
		event.Name,
	)
	if err == sql.ErrNoRows {
		_, insertErr := transaction.NamedExec(`
			INSERT INTO events
			(name, date, start_time, end_time, address, n_of_slots, n_of_volunteers, total_hours, leaders, made_by, sign_up_url)
			VALUES
			(:name, :date, :start_time, :end_time, :address, :n_of_slots, :n_of_volunteers, :total_hours, :leaders, :made_by, :sign_up_url)`,
			event,
		)
		if insertErr != nil {
			return fmt.Errorf("Issue inserting event during upsert: %v", insertErr)
		}
	} else if err != nil {
		return fmt.Errorf("Issue upserting event: %v", err)
	} else {
		event.ID = result.ID // to update the correct row based on primary key (id)
		_, updateErr := transaction.NamedExec(`
			UPDATE events SET 
			name=:name, date=:date, start_time=:start_time, end_time=:end_time, address=:address, n_of_slots=:n_of_slots, n_of_volunteers=:n_of_volunteers, total_hours=:total_hours, leaders=:leaders, made_by=:made_by, sign_up_url=:sign_up_url
			WHERE id=:id
		`, event)
		if updateErr != nil {
			return fmt.Errorf("Issue updating event during upsert: %v", updateErr)
		}
	}
	return nil
}
