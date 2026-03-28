package eventutils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"google.golang.org/api/docs/v1"
)

// logs the event by adding it to the google calendar
// check if the sign up sheet has hours calculated
// if no, calculate them then move on, if yes, move on
// finds first next empty col in the hours sheet and adds it on there
// returns Event struct and an array of members with hours logged and not logged
func LogEvent(signUpDocUrl string, docsService *docs.Service) (LogEventResponse, error) {
	// first, fetch the tables
	// get the metadata table and extract event info from it
	// get the other tables and check if the hours have been calculated
	// if no, calculate them and write them, if yes, move on
	// find next empty column in the EventsMembers sheet and write to it
	// find next empty row in Events and write to it

	documentId := docsUrlToId(signUpDocUrl)

	tables, err := fetchTables(documentId, docsService)
	if err != nil {
		return LogEventResponse{}, fmt.Errorf("Issue logging event: %w", err)
	}

	infoTable := tables[0]
	// attendanceTables := tables[1:]
	getEventInfo(infoTable)

	return LogEventResponse{}, nil
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

// returns an Event object
// doesn't save to db because updates to the db should only happen when fetching data from the sheets
func getEventInfo(tables []docs.Table) (Event, []MemberAttendance) {
	name := getCellText(infoTable.TableRows[0].TableCells[1])
	date := getCellText(infoTable.TableRows[1].TableCells[1])
	times := getCellText(infoTable.TableRows[2].TableCells[1])
	address := getCellText(infoTable.TableRows[3].TableCells[1])

	fmt.Println(name)
	fmt.Println(date)
	fmt.Println(times)
	fmt.Println(address)
	// fmt.Println(infoTable.TableRows[3].TableCells[1].Content[0].Paragraph.Elements[0].TextRun.Content)
	// fmt.Println(infoTable.TableRows[4].TableCells[1].Content[0].Paragraph.Elements[0].TextRun.Content)
	// fmt.Println(infoTable.TableRows[5].TableCells[1].Content[0].Paragraph.Elements[0].TextRun.Content)

	return Event{}
}

// returns member attendance with calculated hours
// ignores grade and phone number because only the name and hours are needed for logging
func getEventAttendance(otherTables *[]docs.Table) ([]MemberAttendance, int) {
	memberAttendance := []MemberAttendance{}
	nOfSlots := 0

	for _, table := range *otherTables {
		// skips the first row because it's a header row
		for _, row := range table.TableRows[1:] {
			// get the name, start and end times
			nOfSlots++
			name := getCellText(row.TableCells[1])
			startTime := getCellText(row.TableCells[4])
			endTime := getCellText(row.TableCells[5])
			memberAttendance = append(memberAttendance, MemberAttendance{
				Name:  name,
				Hours: calculateHours(startTime, endTime),
			})
		}
	}

	return memberAttendance, nOfSlots
}

// returns the contents of a google docs table cell
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
func calculateHours(startTime string, endTime string) float64 {
	startHour, startMins := strings.Split(startTime, ":")
	endHour, endMins := strings.Split(endTime, ":")

	startHourFloat := strconv.ParseFloat(startHour)
	startMinsFloat := strconv.ParseFloat(startMins)
	endHourFloat := strconv.ParseFloat(endHour)
	endMinsFloat := strconv.ParseFloat(endMins)

	// for events that cross noon
	// 9:00 to 3:00, turn to 9:00 to 15:00
	if startHourFloat > endHourFloat {
		endHourFloat += 12
	}

	// calculates hours through minutes and rounds to 2 decimal places
	var elapsedMins float64 = (endHourFloat*60 + endMinsFloat) - (startHourFloat*60 + startMinsFloat)
	var hours float64 = elapsedMins / 60
	hours = math.Round(hours*100) / 100

	return hours
}
