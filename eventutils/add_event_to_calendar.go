package eventutils

import (
	"context"
	"fmt"
	"regexp"

	"keyclubDiscordBot/internal"

	"google.golang.org/api/calendar/v3"
)

// Takes an attendance document ID and adds the event to the calendar, returning the link to the calendar event
func AddEventToCalendar(ctx context.Context, app *internal.App, documentId string) (calendar.Event, error) {
	event, _, err := GetEventInfo(ctx, documentId, app.GoogleServices.Docs)
	if err != nil {
		return calendar.Event{}, fmt.Errorf("Issue extracting event info while adding event to calendar: %w", err)
	}
	noDate := regexp.MustCompile(`\(.*?\)`)
	fmt.Printf("%sT%s-07:00\n", event.Date, event.StartTime)
	fmt.Printf("%sT%s-07:00\n", event.Date, event.EndTime)

	calendarEvent := &calendar.Event{
		Summary:     noDate.ReplaceAllString(event.Name, ""),
		Location:    event.Address,
		Description: "Open the attendance document to view the description.",
		Start: &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT%s-07:00", event.Date, event.StartTime),
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT%s-07:00", event.Date, event.EndTime),
			TimeZone: "America/Los_Angeles",
		},
		Attachments: []*calendar.EventAttachment{
			{
				FileUrl: fmt.Sprintf("https://docs.google.com/document/d/%s/edit?tab=t.0", documentId),
				Title:   "Attendance Document",
			},
		},
	}

	if alreadyExists(ctx, app, calendarEvent) {
		return calendar.Event{}, fmt.Errorf("Event already exists in calendar")
	}

	result, err := app.GoogleServices.Calendar.Events.Insert(app.Config.CalendarID, calendarEvent).Context(ctx).SupportsAttachments(true).Do()
	if err != nil {
		return calendar.Event{}, fmt.Errorf("Issue inserting event into calendar: %w", err)
	}
	return *result, nil
}

func alreadyExists(ctx context.Context, app *internal.App, event *calendar.Event) bool {
	result, err := app.GoogleServices.Calendar.Events.List(app.Config.CalendarID).TimeMin(event.Start.DateTime).TimeMax(event.End.DateTime).Context(ctx).Do()
	if err != nil {
		fmt.Printf("Issue checking if event already exists: %v\n", err)
	}

	for _, item := range result.Items {
		if item.Summary == event.Summary {
			fmt.Printf("Event already exists: %s\n", item.HtmlLink)
			return true
		}
	}
	return false
}
