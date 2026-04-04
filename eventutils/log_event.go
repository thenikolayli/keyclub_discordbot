package eventutils

import (
	"fmt"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/memberutils"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"google.golang.org/api/docs/v1"
)

func LogEvent(signUpDocUrl string) (LogEventResponse, error) {
	documentId := docsUrlToId(signUpDocUrl)
	event, memberAttendance, err := getEventInfo(documentId, config.GoogleServices.Docs)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue extracting event info while logging event: %w", err)
	}

	requests := batchRequests(memberAttendance)
	config.GoogleServices.Docs.Documents.BatchUpdate(config.SpreadsheetID, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()

	emptyRow, err := findNextEmptyRowNoDupes(config.EventsMembersSheetRanges.Events, event.Name)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue finding next empty row while logging event: %w", err)
	}
	// make a function that gets member columns

	fmt.Println(event, emptyRow)
	findMemberColumns(&memberAttendance)

	for _, member := range memberAttendance {
		fmt.Println(member.Name, member.Column)
	}
	return LogEventResponse{}, nil
}

// creates a batch of updates to the sign up doc (to write calculated hours)
func batchRequests(memberAttendance []MemberAttendance) []*docs.Request {
	requests := []*docs.Request{}
	for _, member := range memberAttendance {
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
func findNextEmptyRowNoDupes(searchRange string, eventName string) (string, error) {
	response, err := config.GoogleServices.Sheets.Spreadsheets.Values.Get(config.SpreadsheetID, searchRange).Do()
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

// finds the columns of members in the EventsMembers sheet
func findMemberColumns(memberAttendance *[]MemberAttendance) error {
	// response, err := config.GoogleServices.Sheets.Spreadsheets.Values.Get(config.SpreadsheetID, config.EventsMembersSheetRanges.Members).Do()
	namesResponse, err := config.GoogleServices.Sheets.Spreadsheets.Values.Get(config.SpreadsheetID, config.EventsMembersSheetRanges.Members).Do()
	if err != nil {
		return fmt.Errorf("Issue fetching member columns from sheet while logging event: %w", err)
	}
	nicknamesResponse, err := config.GoogleServices.Sheets.Spreadsheets.Values.Get(config.SpreadsheetID, config.EventsMembersSheetRanges.MemberNicknames).Do()
	if err != nil {
		return fmt.Errorf("Issue fetching member nicknames from sheet while logging event: %w", err)
	}

	for index, memberValue := range namesResponse.Values[0] {
		for memberIndex := range *memberAttendance {
			if (*memberAttendance)[memberIndex].Column == "" {
				// compare member value name to member attendance name
				memberValueName := memberutils.NewName(memberValue.(string))
				memberAttendanceValueName := memberutils.NewName((*memberAttendance)[memberIndex].Name)
				fmt.Printf("comparing sheet: %+v | attendance: %+v\n", memberValueName, memberAttendanceValueName)

				if memberAttendanceValueName.Firstname == memberValueName.Firstname && memberAttendanceValueName.Lastname == memberValueName.Lastname {
					(*memberAttendance)[memberIndex].Column = indexToCol(index)
				} else if memberAttendanceValueName.Firstname == memberValueName.Lastname && memberAttendanceValueName.Lastname == memberValueName.Firstname {
					(*memberAttendance)[memberIndex].Column = indexToCol(index)
				} else if memberAttendanceValueName.Nickname != "" && memberAttendanceValueName.Nickname == memberValueName.Nickname || memberAttendanceValueName.Firstname == memberValueName.Firstname {
					(*memberAttendance)[memberIndex].Column = indexToCol(index)
				}
			}
		}
	}

	return nil
}

// converts a numerical index to a column letter (1 -> A, 2 -> B, 27 -> AA, etc)
func indexToCol(index int) string {
	result := ""
	for index > 0 {
		index-- // shift to 0-indexed
		result = string(rune('A'+index%26)) + result
		index /= 26
	}
	return result
}

// gets the event info to put in the Events sheet (excluding tags, since that has to be done manually)
// it's just the event info from the sign up sheet, so whatever is on there gets saved
// doesn't save to db because updates to the db should only happen when fetching data from the sheets
func getEventInfo(documentId string, docsService *docs.Service) (Event, []MemberAttendance, error) {
	tables, err := fetchTables(documentId, docsService)
	if err != nil {
		return Event{}, []MemberAttendance{}, fmt.Errorf("Error fetching tables: %w", err)
	}

	infoTable := tables[0]
	attendanceTables := tables[1:]

	name := getCellText(infoTable.TableRows[0].TableCells[1])
	address := getCellText(infoTable.TableRows[3].TableCells[1])
	dateString := getCellText(infoTable.TableRows[1].TableCells[1])
	timeString := getCellText(infoTable.TableRows[2].TableCells[1])
	date, startTime, endTime := parseDateAndTime(dateString, timeString)

	nOfSlots := 0
	nOfVolunteers := 0
	totalHours := 0.0
	memberAttendance := []MemberAttendance{}

	for _, table := range attendanceTables {
		// skips the header row
		for _, row := range table.TableRows[1:] {
			nOfSlots++
			name := getCellText(row.TableCells[1])
			if name != "" {
				nOfVolunteers++
				startTime := getCellText(row.TableCells[4])
				endTime := getCellText(row.TableCells[5])
				hours, err := calculateHours(startTime, endTime)
				if err != nil {
					return Event{}, []MemberAttendance{}, fmt.Errorf("error calculating hours for member %s: %w", name, err)
				}
				totalHours += hours
				memberAttendance = append(memberAttendance, MemberAttendance{
					Name:  name,
					Hours: hours,
					// + 1 and - 1 offsets so it grabs the inside of the cell, rather than the borders
					HoursStartIndex: int(row.TableCells[6].StartIndex) + 1,
					HoursEndIndex:   int(row.TableCells[6].EndIndex) - 1,
				})
			}
		}
	}

	return Event{
		Name:          name,
		Date:          date,
		StartTime:     startTime,
		EndTime:       endTime,
		Address:       address,
		NofSlots:      nOfSlots,
		NofVolunteers: nOfVolunteers,
		TotalHours:    totalHours,
	}, memberAttendance, nil
}

// returns the text contents of a google docs table cell
func getCellText(tableCell *docs.TableCell) string {
	var stringBuilder strings.Builder
	for _, content := range tableCell.Content {
		if content.Paragraph == nil {
			continue
		}
		for _, elem := range content.Paragraph.Elements {
			if elem.TextRun != nil {
				stringBuilder.WriteString(elem.TextRun.Content)
			}
		}
	}
	return strings.TrimSpace(stringBuilder.String())
}

// calculates hours between start and end time
func calculateHours(startTime string, endTime string) (float64, error) {
	start := strings.Split(startTime, ":")
	end := strings.Split(endTime, ":")

	// ignoring error isn't ideal, but i'll fix this
	startHourFloat, err := strconv.ParseFloat(start[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing start hour (possibly empty cells): %w", err)
	}
	startMinsFloat, err := strconv.ParseFloat(start[1], 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing start minutes (possibly empty cells): %w", err)
	}
	endHourFloat, err := strconv.ParseFloat(end[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing end hour (possibly empty cells): %w", err)
	}
	endMinsFloat, err := strconv.ParseFloat(end[1], 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing end minutes (possibly empty cells): %w", err)
	}

	// for events that cross noon
	// 9:00 to 3:00, turn to 9:00 to 15:00
	if startHourFloat > endHourFloat {
		endHourFloat += 12
	}

	// calculates hours through minutes and rounds to 2 decimal places
	var elapsedMins float64 = (endHourFloat*60 + endMinsFloat) - (startHourFloat*60 + startMinsFloat)
	var hours float64 = elapsedMins / 60
	hours = math.Round(hours*100) / 100

	return hours, nil
}

// parses date and start and end time strings into time.Time objects
// ignoring errors isn't ideal but i'll fix this later...
func parseDateAndTime(dateString string, timeString string) (string, string, string) {
	date, _ := dateparse.ParseAny(dateString)
	dateFormatted := date.Format(time.DateOnly)

	timeParts := strings.Split(timeString, "-")
	if len(timeParts) != 2 {
		timeParts = strings.Split(timeString, "to")
	}
	// adding a date so dateparse can parse the time, then just take the date away
	startTime, _ := dateparse.ParseAny(fmt.Sprintf("January 1, 2000 %v", timeParts[0]))
	startTimeFormatted := startTime.Format(time.TimeOnly)
	endTime, _ := dateparse.ParseAny(fmt.Sprintf("January 1, 2000 %v", timeParts[1]))
	endTimeFormatted := endTime.Format(time.TimeOnly)

	return dateFormatted, startTimeFormatted, endTimeFormatted
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

// returns all the tables of the document
func fetchTables(documentId string, docsService *docs.Service) ([]docs.Table, error) {
	document, err := docsService.Documents.Get(documentId).Do()
	tables := []docs.Table{}
	if err != nil {
		return []docs.Table{}, fmt.Errorf("Issue fetching document: %w", err)
	}

	for _, structuralElement := range document.Body.Content {
		if structuralElement.Table != nil {
			tables = append(tables, *structuralElement.Table)
		}
	}

	return tables, nil
}

// https://docs.google.com/document/d/id-example/edit?tab=t.0
// extracts a Google docs id from the url
// splits it at /d/ and gets the second part, then splits it at /edit and gets the first one
func docsUrlToId(url string) string {
	id := strings.Split(url, "/d/")[1]
	id = strings.Split(id, "/edit")[0]
	return id
}
