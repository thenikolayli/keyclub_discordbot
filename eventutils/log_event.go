package eventutils

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"keyclubDiscordBot/internal"
	"keyclubDiscordBot/memberutils"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/sheets/v4"
)

// Takes an attendance document ID and fills calculated hours and logs event to the sheets, returning the event info and which members were logged vs not logged
func LogEvent(ctx context.Context, app *internal.App, documentId string) (LogEventResponse, error) {
	event, memberAttendance, err := GetEventInfo(ctx, documentId, app.GoogleServices.Docs)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue extracting event info while logging event: %w", err)
	}
	requests := batchRequests(memberAttendance)
	_, err = app.GoogleServices.Docs.Documents.BatchUpdate(documentId, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Context(ctx).Do()
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue updating document while logging event: %w", err)
	}

	emptyRowEventsMembers, err := findNextEmptyRowNoDupes(ctx, app, app.Config.EventsMembersSheetRanges.Events, event.Name)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue finding next empty row while logging event: %w", err)
	}
	eventsMembersUpdateValues, logEventResponse, err := createUpdateValues(ctx, app, &memberAttendance, event)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue creating update values while logging event: %w", err)
	}
	emptyRowEvents, err := findNextEmptyRowNoDupes(ctx, app, app.Config.EventsSheetRanges.Events, event.Name)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue finding next empty row while logging event: %w", err)
	}
	eventsUpdateValues := []any{
		event.Name,
		event.Date,
		event.StartTime,
		event.EndTime,
		event.Address,
		event.NofSlots,
		event.NofVolunteers,
		event.TotalHours,
		event.Leaders,
		event.MadeBy,
	}

	_, err = app.GoogleServices.Sheets.Spreadsheets.Values.BatchUpdate(
		app.Config.SpreadsheetID,
		&sheets.BatchUpdateValuesRequest{
			ValueInputOption: "USER_ENTERED",
			Data: []*sheets.ValueRange{
				{
					Range:  fmt.Sprintf("%s!A%v:%s%v", app.Config.EventsMembersSheetRanges.SheetName, emptyRowEventsMembers, indexToCol(len(eventsMembersUpdateValues)), emptyRowEventsMembers),
					Values: [][]any{eventsMembersUpdateValues},
				},
				{ // goes to H because tags aren't logged (they have to be added manually)
					Range:  fmt.Sprintf("%s!A%v:J%v", app.Config.EventsSheetRanges.SheetName, emptyRowEvents, emptyRowEvents),
					Values: [][]any{eventsUpdateValues},
				},
			},
		},
	).Context(ctx).Do()
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue updating sheets while logging event: %w", err)
	}

	return logEventResponse, nil
}

// finds the columns of members in the EventsMembers sheet
func createUpdateValues(ctx context.Context, app *internal.App, memberAttendance *[]MemberAttendance, event Event) ([]any, LogEventResponse, error) {
	namesResponse, err := app.GoogleServices.Sheets.Spreadsheets.Values.Get(app.Config.SpreadsheetID, app.Config.EventsMembersSheetRanges.Members).Context(ctx).Do()
	if err != nil {
		return nil, LogEventResponse{}, fmt.Errorf("Issue fetching member columns from sheet while logging event: %w", err)
	}
	updateValues := make([]any, len(namesResponse.Values[0]))
	membersLogged := []MemberAttendance{}
	membersNotLogged := []MemberAttendance{}

	for index, sheetValue := range namesResponse.Values[0] {
		sheetName := memberutils.NewName(sheetValue.(string))
		for memberIndex := range *memberAttendance {
			if (*memberAttendance)[memberIndex].ColumnFound || updateValues[index] != nil {
				continue
			}
			memberAttendanceName := memberutils.NewName((*memberAttendance)[memberIndex].Name)
			if memberutils.SameName(memberAttendanceName, sheetName) {
				if (*memberAttendance)[memberIndex].Hours != -1 {
					updateValues[index] = (*memberAttendance)[memberIndex].Hours
				} else {
					updateValues[index] = nil
				}
				(*memberAttendance)[memberIndex].ColumnFound = true
			} else {
				updateValues[index] = nil
			}
		}
	}

	for _, member := range *memberAttendance {
		if member.ColumnFound && member.Hours != -1 {
			membersLogged = append(membersLogged, member)
		} else {
			membersNotLogged = append(membersNotLogged, member)
		}
	}
	updateValues[0] = event.Name
	return updateValues, LogEventResponse{
		Event:            event,
		MembersLogged:    membersLogged,
		MembersNotLogged: membersNotLogged,
	}, nil
}

// writes/overwrites calculated hours
func writeHoursToCell(memberAttendance MemberAttendance) (*docs.DeleteContentRangeRequest, *docs.InsertTextRequest) {
	// if the cell is blank, just insert it
	// 1  2  3 vs 1  2 3 4  5
	// 1 3 5 are the borders and 2 3 4 are inside,
	// if start + 1 == end - 1 then it's blank, otherwise it has content
	if memberAttendance.HoursEndIndex == memberAttendance.HoursStartIndex {
		return nil, &docs.InsertTextRequest{
			Text: strconv.FormatFloat(memberAttendance.Hours, 'f', 2, 64),
			Location: &docs.Location{
				Index: int64(memberAttendance.HoursStartIndex),
			},
		}
	}
	// if it's not blank, delete the existing content and insert the new hours
	return &docs.DeleteContentRangeRequest{
			Range: &docs.Range{
				StartIndex: int64(memberAttendance.HoursStartIndex),
				EndIndex:   int64(memberAttendance.HoursEndIndex),
			},
		}, &docs.InsertTextRequest{
			Text: strconv.FormatFloat(memberAttendance.Hours, 'f', 2, 64),
			Location: &docs.Location{
				Index: int64(memberAttendance.HoursStartIndex),
			},
		}
}

// creates a batch of updates to the sign up doc (to write calculated hours)
func batchRequests(memberAttendance []MemberAttendance) []*docs.Request {
	requests := []*docs.Request{}
	for _, member := range memberAttendance {
		// Treat -1 hours as "not logged" (leave blank in doc).
		if member.Hours == -1 {
			continue
		}
		deleteRequest, insertRequest := writeHoursToCell(member)
		requests = append(requests, &docs.Request{
			InsertText: insertRequest,
		})
		// only check if deleterequest is nil because insert text request will always not be nil
		// delete requests are appended after insert because the list is reversed
		if deleteRequest != nil {
			requests = append(requests, &docs.Request{
				DeleteContentRange: deleteRequest,
			})
		}
	}
	// reverses so updates happen backwards and don't affect/offset start indexes
	slices.Reverse(requests)
	return requests
}

// finds next empty row of sheet and makes sure the event isn't already logged
// takes searchRange because it can be used for Events and EventsMembers sheet
func findNextEmptyRowNoDupes(ctx context.Context, app *internal.App, searchRange string, eventName string) (string, error) {
	response, err := app.GoogleServices.Sheets.Spreadsheets.Values.Get(app.Config.SpreadsheetID, searchRange).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("Issue fetching events from sheet while logging event: %w", err)
	}
	for _, row := range response.Values {
		if row[0] == eventName {
			return "", fmt.Errorf("Event %v already logged in sheet", eventName)
		}
	}
	return strconv.Itoa(len(response.Values) + 1), nil
}

// converts a numerical index to a column letter (1 -> A, 2 -> B, 27 -> AA, etc)
func indexToCol(index int) string {
	result := ""
	index++
	for index > 0 {
		index-- // shift to 0-indexed
		result = string(rune('A'+index%26)) + result
		index /= 26
	}
	return result
}
