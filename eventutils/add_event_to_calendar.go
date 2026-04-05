package eventutils

import (
	"fmt"
	"regexp"

	"keyclubDiscordBot/config"

	"google.golang.org/api/calendar/v3"
)

// Takes an attendance document ID and adds the event to the calendar, returning the link to the calendar event
func AddEventToCalendar(documentId string) (string, error) {
	event, _, err := GetEventInfo(documentId, config.GoogleServices.Docs)
	if err != nil {
		return "", fmt.Errorf("Issue extracting event info while adding event to calendar: %w", err)
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

	if alreadyExists(calendarEvent) {
		return "", fmt.Errorf("Event already exists in calendar")
	}

	result, err := config.GoogleServices.Calendar.Events.Insert(config.CalendarID, calendarEvent).SupportsAttachments(true).Do()
	if err != nil {
		return "", fmt.Errorf("Issue inserting event into calendar: %w", err)
	}
	return result.HtmlLink, nil
}

func alreadyExists(event *calendar.Event) bool {
	result, err := config.GoogleServices.Calendar.Events.List(config.CalendarID).TimeMin(event.Start.DateTime).TimeMax(event.End.DateTime).Do()
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
