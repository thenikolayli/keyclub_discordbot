package eventutils

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"google.golang.org/api/docs/v1"
)

// gets the event info to put in the Events sheet (excluding tags, since that has to be done manually)
// it's just the event info from the sign up sheet, so whatever is on there gets saved
// doesn't save to db because updates to the db should only happen when fetching data from the sheets
func GetEventInfo(ctx context.Context, documentId string, docsService *docs.Service) (Event, []MemberAttendance, error) {
	tables, err := fetchTables(ctx, documentId, docsService)
	if err != nil {
		return Event{}, []MemberAttendance{}, fmt.Errorf("Error fetching tables: %w", err)
	}

	infoTable := tables[0]
	attendanceTables := tables[1:]

	name := ""
	address := ""
	dateString := ""
	timeString := ""
	madeBy := ""
	leaders := ""

	for _, row := range infoTable.TableRows {
		header := strings.ToLower(getCellText(row.TableCells[0]))
		value := getCellText(row.TableCells[1])
		switch header {
		case "event name:":
			name = value
		case "event:":
			name = value
		case "address:":
			address = value
		case "leaders:":
			leaders = value
		case "made by:":
			madeBy = value
		case "date:":
			dateString = value
		case "time:":
			timeString = value
		}
	}
	date, nameDate, startTime, endTime, err := parseDateAndTime(dateString, timeString)
	if err != nil {
		return Event{}, []MemberAttendance{}, fmt.Errorf("error parsing event date/time: %w", err)
	}

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
				if hours != -1 {
					totalHours += hours
				}
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

	return Event{ // (m/d) name
		Name:          fmt.Sprintf("%s %s", nameDate, name),
		Date:          date,
		StartTime:     startTime,
		EndTime:       endTime,
		Address:       address,
		NofSlots:      nOfSlots,
		NofVolunteers: nOfVolunteers,
		TotalHours:    totalHours,
		Leaders:       leaders,
		MadeBy:        madeBy,
		SignUpUrl:     fmt.Sprintf("https://docs.google.com/document/d/%s/edit?tab=t.0", documentId),
	}, memberAttendance, nil
}

// calculates hours between start and end time
func calculateHours(startTime string, endTime string) (float64, error) {
	startTime = strings.TrimSpace(startTime)
	endTime = strings.TrimSpace(endTime)
	if startTime == "" && endTime == "" {
		return -1, nil
	}

	start := strings.Split(startTime, ":")
	end := strings.Split(endTime, ":")
	if len(start) != 2 || len(end) != 2 {
		return 0, fmt.Errorf("invalid time format (expected HH:MM): start=%q end=%q", startTime, endTime)
	}

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
func fetchTables(ctx context.Context, documentId string, docsService *docs.Service) ([]docs.Table, error) {
	document, err := docsService.Documents.Get(documentId).Context(ctx).Do()
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
func parseDateAndTime(dateString string, timeString string) (string, string, string, string, error) {
	normalizedDateString := normalizeDateString(dateString)
	date, err := dateparse.ParseAny(normalizedDateString)
	if err != nil {
		return "", "", "", "", fmt.Errorf("could not parse date %q (normalized: %q): %w", dateString, normalizedDateString, err)
	}
	dateFormatted := date.Format(time.DateOnly)
	nameDate := fmt.Sprintf("(%d/%d)", int(date.Month()), date.Day())

	timeParts := strings.Split(timeString, "-")
	if len(timeParts) != 2 {
		timeParts = strings.Split(timeString, "to")
	}
	if len(timeParts) != 2 {
		return "", "", "", "", fmt.Errorf("could not split time range %q", timeString)
	}
	// adding a date so dateparse can parse the time, then just take the date away
	startTime, err := dateparse.ParseAny(fmt.Sprintf("January 1, 2000 %v", timeParts[0]))
	if err != nil {
		return "", "", "", "", fmt.Errorf("could not parse start time %q: %w", strings.TrimSpace(timeParts[0]), err)
	}
	startTimeFormatted := startTime.Format(time.TimeOnly)
	endTime, err := dateparse.ParseAny(fmt.Sprintf("January 1, 2000 %v", timeParts[1]))
	if err != nil {
		return "", "", "", "", fmt.Errorf("could not parse end time %q: %w", strings.TrimSpace(timeParts[1]), err)
	}
	endTimeFormatted := endTime.Format(time.TimeOnly)

	return dateFormatted, nameDate, startTimeFormatted, endTimeFormatted, nil
}

var (
	weekdayPrefixRegex = regexp.MustCompile(`(?i)^\s*(mon(day)?|tue(s(day)?)?|wed(nesday)?|thu(r(s(day)?)?)?|fri(day)?|sat(urday)?|sun(day)?)\s*,?\s*`)
	ordinalRegex       = regexp.MustCompile(`(?i)(\d+)(st|nd|rd|th)\b`)
)

func normalizeDateString(dateString string) string {
	s := strings.TrimSpace(dateString)
	s = strings.ReplaceAll(s, ",", " ")
	s = weekdayPrefixRegex.ReplaceAllString(s, "")
	s = ordinalRegex.ReplaceAllString(s, "$1")
	// collapse whitespace
	return strings.Join(strings.Fields(s), " ")
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
