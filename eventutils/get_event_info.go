package eventutils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"google.golang.org/api/docs/v1"
)

// gets the event info to put in the Events sheet (excluding tags, since that has to be done manually)
// it's just the event info from the sign up sheet, so whatever is on there gets saved
// doesn't save to db because updates to the db should only happen when fetching data from the sheets
func GetEventInfo(documentId string, docsService *docs.Service) (Event, []MemberAttendance, error) {
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

	return Event{ // ISO 8096 format (yyyy-mm-dd) name
		Name:          fmt.Sprintf("(%s) %s", date, name),
		Date:          date,
		StartTime:     startTime,
		EndTime:       endTime,
		Address:       address,
		NofSlots:      nOfSlots,
		NofVolunteers: nOfVolunteers,
		TotalHours:    totalHours,
	}, memberAttendance, nil
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
