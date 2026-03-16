package hoursutils

import (
	"fmt"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/genericutils"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// when a user requests hours, the program should check if their hours have been

// func GetHours(name string) *Hours {
// 	now := time.Now().Unix()

// 	if lastUpdated < now-config.HoursTTL {

// 	}
// }

// steps
// get values from google sheets api
// format and turn them all into structs

// updates the member database entries
// fetches values via an api call to the hours spreadsheet
// formats the response to member structs
// updates the database based on structs
func UpdateMembers(googleServices *genericutils.GoogleServices, db *sqlx.DB) error {
	prevTime := time.Now()
	memberValueRanges, err := getMemberValueRanges(googleServices)
	if err != nil {
		return fmt.Errorf("Failed to update members: %v", err)
	}
	fmt.Printf("Time for API call: %v\n", time.Since(prevTime))

	prevTime = time.Now()
	formattedMemberStructs := getFormattedMemberStructs(memberValueRanges)
	fmt.Printf("Time to format: %v\n", time.Since(prevTime))

	prevTime = time.Now()
	// check if it exists first, if yes, update, if no, add it
	for _, each := range formattedMemberStructs {
		err := upsertMember(each, db)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Time to run DB queries: %v", time.Since(prevTime))

	return nil
}

// upserts member
func upsertMember(member Member, db *sqlx.DB) error {
	_, err := db.NamedExec(`
		insert into members
		(name, nickname, term_hours, all_hours, shirt_size, paid_dues, grad_year, strikes)
		values
		(:name, :nickname, :term_hours, :all_hours, :shirt_size, :paid_dues, :grad_year, :strikes)
		on conflict(name) do update set
		nickname=excluded.nickname, 
		term_hours=excluded.term_hours, 
		all_hours=excluded.all_hours, 
		shirt_size=excluded.shirt_size, 
		paid_dues=excluded.paid_dues, 
		grad_year=excluded.grad_year, 
		strikes=excluded.strikes
	`, member)
	if err != nil {
		return fmt.Errorf("Error upserting %v: %v\n", member.Name, err)
	}
	return nil
}

// fetches and returns google sheets api value ranges (unformatted)
func getMemberValueRanges(googleServices *genericutils.GoogleServices) ([]*sheets.ValueRange, error) {
	data, err := googleServices.Sheets.Spreadsheets.Values.BatchGet(config.SpreadsheetID).Ranges(
		config.NamesRange,
		config.NicknamesRange,
		config.TermHoursRange,
		config.AllHoursRange,
		config.GradYearRange,
		config.StrikesRange,
	).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to batch get spreadsheet ranges: %v", err)
	}
	return data.ValueRanges, nil
}

// takes the api call value ranges and turns them into an array of member structs
func getFormattedMemberStructs(memberValueRanges []*sheets.ValueRange) []Member {
	// gets length based on the length of the names column
	memberValueRangesLength := len(memberValueRanges[0].Values)
	formattedMemberArray := make([]Member, memberValueRangesLength)

	normalizedNames := normalizeStringValues(memberValueRanges[0].Values, memberValueRangesLength)
	normalizedNicknames := normalizeStringValues(memberValueRanges[1].Values, memberValueRangesLength)
	normalizedTermHours := normalizeFloatValues(memberValueRanges[2].Values, memberValueRangesLength)
	normalizedAllHours := normalizeFloatValues(memberValueRanges[3].Values, memberValueRangesLength)
	normalizedGradYears := normalizeIntValues(memberValueRanges[4].Values, memberValueRangesLength)
	normalizedStrikes := normalizeIntValues(memberValueRanges[5].Values, memberValueRangesLength)

	for i := range memberValueRangesLength - 1 {
		formattedMemberArray[i] = Member{
			Name:      normalizedNames[i],
			Nickname:  normalizedNicknames[i],
			TermHours: normalizedTermHours[i],
			AllHours:  normalizedAllHours[i],
			GradYear:  normalizedGradYears[i],
			Strikes:   normalizedStrikes[i],
			ShirtSize: "",
			PaidDues:  false,
			ID:        -1,
		}
	}

	return formattedMemberArray
}

// the following functions create normalized lists with blanks from the originals
func normalizeStringValues(values [][]any, length int) []string {
	normalizedStringValues := make([]string, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			normalizedStringValues[i] = values[i][0].(string)
		}
	}

	return normalizedStringValues
}

func normalizeFloatValues(values [][]any, length int) []float64 {
	normalizedStringValues := make([]float64, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			value, err := strconv.ParseFloat(values[i][0].(string), 64)
			if err == nil {
				normalizedStringValues[i] = value
			}
		}
	}

	return normalizedStringValues
}

func normalizeIntValues(values [][]any, length int) []int {
	normalizedStringValues := make([]int, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			value, err := strconv.ParseInt(values[i][0].(string), 0, 64)
			if err == nil {
				normalizedStringValues[i] = int(value)
			}
		}
	}

	return normalizedStringValues
}
