package eventutils

import (
	"database/sql"
	"fmt"
	"time"

	"keyclubDiscordBot/config"
	"keyclubDiscordBot/genericutils"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

func SyncEvents(nSlots int, eventsUpdateTimeout float64, eventsLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) error {
	return syncEventsFromSheet(eventsUpdateTimeout, eventsLastUpdated, sheetsService, database)
}

// syncs the events database entries from the events spreadsheet
// fetches values via an api call to the events spreadsheet
// formats the response to event structs
// updates the database based on structs
func syncEventsFromSheet(eventsUpdateTimeout float64, eventsLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) error {
	timeSince := time.Since(*eventsLastUpdated).Seconds()
	if timeSince < eventsUpdateTimeout {
		return fmt.Errorf("Not enough time has passed since the last update, wait %v more seconds.", eventsUpdateTimeout-timeSince)
	}

	eventValueRanges, err := getEventValueRanges(sheetsService)
	if err != nil {
		return fmt.Errorf("Failed to update events: %v", err)
	}

	formattedEventStructs := getEventStructs(eventValueRanges)

	transaction, err := database.BeginTxx(config.Context, nil)
	if err != nil {
		return fmt.Errorf("Failed to create a transaction: %v", err)
	}
	for _, each := range formattedEventStructs {
		err := upsertEvent(each, transaction)
		if err != nil {
			return err
		}
	}
	transaction.Commit()

	*eventsLastUpdated = time.Now()
	return nil
}

// fetches and returns google sheets api value ranges (unformatted)
func getEventValueRanges(sheetsService *sheets.Service) ([]*sheets.ValueRange, error) {
	data, err := sheetsService.Spreadsheets.Values.BatchGet(config.SpreadsheetID).Ranges(
		config.EventsSheetRanges.Events,
		config.EventsSheetRanges.Dates,
		config.EventsSheetRanges.StartTimes,
		config.EventsSheetRanges.EndTimes,
		config.EventsSheetRanges.Addresses,
		config.EventsSheetRanges.NofSlots,
		config.EventsSheetRanges.NofVolunteers,
		config.EventsSheetRanges.TotalHours,
		config.EventsSheetRanges.Leaders,
		config.EventsSheetRanges.MadeBy,
	).Do()
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

	for i := range eventValueRangesLength - 1 {
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
func upsertEvent(event Event, transaction *sqlx.Tx) error {
	result := Event{}
	err := transaction.GetContext(
		config.Context, &result,
		"SELECT * from events WHERE name = ? LIMIT 1",
		event.Name,
	)
	if err == sql.ErrNoRows {
		_, insertErr := transaction.NamedExec(`
			INSERT INTO events
			(name, date, start_time, end_time, address, n_of_slots, n_of_volunteers, total_hours, leaders, made_by)
			VALUES
			(:name, :date, :start_time, :end_time, :address, :n_of_slots, :n_of_volunteers, :total_hours, :leaders, :made_by)`,
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
			name=:name, date=:date, start_time=:start_time, end_time=:end_time, address=:address, n_of_slots=:n_of_slots, n_of_volunteers=:n_of_volunteers, total_hours=:total_hours, leaders=:leaders, made_by=:made_by
			WHERE id=:id
		`, event)
		if updateErr != nil {
			return fmt.Errorf("Issue updating event during upsert: %v", updateErr)
		}
	}
	return nil
}
