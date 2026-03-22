package memberutils

import (
	"database/sql"
	"fmt"
	"keyclubDiscordBot/config"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// updates the member database entries
// fetches values via an api call to the hours spreadsheet
// formats the response to member structs
// updates the database based on structs
func UpdateMembers(hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) error {
	timeSince := time.Since(*hoursLastUpdated).Seconds()
	if timeSince < hoursUpdateTimeout {
		return fmt.Errorf("Not enough time has passed since the last update, wait %v more seconds.", hoursUpdateTimeout-timeSince)
	}

	prevTime := time.Now()
	memberValueRanges, err := getMemberValueRanges(sheetsService)
	if err != nil {
		return fmt.Errorf("Failed to update members: %v", err)
	}
	fmt.Printf("Time for API call: %v\n", time.Since(prevTime))

	prevTime = time.Now()
	formattedMemberStructs := getFormattedMemberStructs(memberValueRanges)
	fmt.Printf("Time to format: %v\n", time.Since(prevTime))

	// check if it exists first, if yes, update, if no, add it
	prevTime = time.Now()
	transaction, err := database.BeginTxx(config.Context, nil)
	if err != nil {
		return fmt.Errorf("Failed to create a transaction: %v", err)
	}
	for _, each := range formattedMemberStructs {
		err := upsertMember(each, transaction)
		if err != nil {
			return err
		}
	}
	transaction.Commit()
	fmt.Printf("Time to run DB queries: %v\n", time.Since(prevTime))

	*hoursLastUpdated = time.Now()
	return nil
}

// takes in a formatted member struct and a transaction and upserts their row
// checks if a member with the same name exists
// if they don't, insert them
// otherwise, update
func upsertMember(member Member, transaction *sqlx.Tx) error {
	result := Member{}
	err := transaction.GetContext(
		config.Context, &result,
		"SELECT * from members WHERE first_name = ? AND last_name = ? LIMIT 1",
		member.Firstname, member.Lastname,
	)
	if err == sql.ErrNoRows {
		_, insertErr := transaction.NamedExec(`
			INSERT INTO members
			(first_name, last_name, nickname, all_hours, term_hours, grad_year, class_year, strikes, personal_email, school_email, phone_number, shirt_size, paid_dues)
			VALUES
			(:first_name, :last_name, :nickname, :all_hours, :term_hours, :grad_year, :class_year, :strikes, :personal_email, :school_email, :phone_number, :shirt_size, :paid_dues)`,
			member,
		)
		if insertErr != nil {
			return fmt.Errorf("Issue inserting member during upsert: %v", insertErr)
		}
	} else if err != nil {
		return fmt.Errorf("Issue upserting member: %v", err)
	} else {
		member.ID = result.ID // to update the correct row based on primary key (id)
		_, updateErr := transaction.NamedExec(`
			UPDATE members SET 
			first_name=:first_name, last_name=:last_name, nickname=:nickname, all_hours=:all_hours, term_hours=:term_hours, grad_year=:grad_year, class_year=:class_year, strikes=:strikes, personal_email=:personal_email, school_email=:school_email, phone_number=:phone_number, shirt_size=:shirt_size, paid_dues=:paid_dues
			WHERE id=:id
		`, member)
		if updateErr != nil {
			return fmt.Errorf("Issue updating member during upsert: %v", updateErr)
		}
	}
	return nil
}

// fetches and returns google sheets api value ranges (unformatted)
func getMemberValueRanges(sheetsService *sheets.Service) ([]*sheets.ValueRange, error) {
	data, err := sheetsService.Spreadsheets.Values.BatchGet(config.SpreadsheetID).Ranges(
		config.NamesRange,
		config.NicknamesRange,
		config.AllHoursRange,
		config.TermHoursRange,
		config.GradYearRange,
		config.ClassYearRange,
		config.StrikesRange,
		config.PersonalEmailRange,
		config.SchoolEmailRange,
		config.PhoneNumberRange,
		config.ShirtSizesRange,
		config.PaidDuesRange,
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
	normalizedAllHours := normalizeFloatValues(memberValueRanges[2].Values, memberValueRangesLength)
	normalizedTermHours := normalizeFloatValues(memberValueRanges[3].Values, memberValueRangesLength)
	normalizedGradYears := normalizeIntValues(memberValueRanges[4].Values, memberValueRangesLength)
	normalizedClassYears := normalizeIntValues(memberValueRanges[5].Values, memberValueRangesLength)
	normalizedStrikes := normalizeIntValues(memberValueRanges[6].Values, memberValueRangesLength)
	normalizedPersonalEmails := normalizeStringValues(memberValueRanges[7].Values, memberValueRangesLength)
	normalizedSchoolEmails := normalizeStringValues(memberValueRanges[8].Values, memberValueRangesLength)
	normalizedPhoneNumbers := normalizeStringValues(memberValueRanges[9].Values, memberValueRangesLength)
	normalizedShirtSizes := normalizeStringValues(memberValueRanges[10].Values, memberValueRangesLength)
	normalizedPaidDues := normalizeBoolValues(memberValueRanges[11].Values, memberValueRangesLength)
	

	for i := range memberValueRangesLength - 1 {
		name := strings.Split(normalizedNames[i], ",") // names are stored as Last, First in the spreadsheet

		formattedMemberArray[i] = Member{
			Firstname: strings.ToLower(strings.TrimSpace(name[1])),
			Lastname:  strings.ToLower(strings.TrimSpace(name[0])),
			Nickname:  strings.ToLower(normalizedNicknames[i]),
			AllHours:  normalizedAllHours[i],
			TermHours: normalizedTermHours[i],
			GradYear:  normalizedGradYears[i],
			ClassYear: normalizedClassYears[i],
			PersonalEmail: normalizedPersonalEmails[i],
			SchoolEmail:   normalizedSchoolEmails[i],
			PhoneNumber:   normalizedPhoneNumbers[i],
			Strikes:   normalizedStrikes[i],
			ShirtSize: normalizedShirtSizes[i],
			PaidDues:  normalizedPaidDues[i],
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

func normalizeBoolValues(values [][]any, length int) []bool {
	normalizedStringValues := make([]bool, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			value, err := strconv.ParseBool(values[i][0].(string))
			if err == nil {
				normalizedStringValues[i] = bool(value)
			}
		}
	}

	return normalizedStringValues
}
