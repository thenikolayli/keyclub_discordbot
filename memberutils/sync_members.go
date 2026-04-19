package memberutils

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/internal"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// updates the member database entries
// fetches values via an api call to the hours spreadsheet
// formats the response to member structs
// updates the database based on structs
func SyncMembersFromSheet(ctx context.Context, app *internal.App) error {
	if !app.MemberSync.ShouldSync() {
		remaining := app.MemberSync.UpdateTimeout - time.Since(app.MemberSync.LastUpdated)
		return fmt.Errorf("Not enough time has passed since the last update, wait %v more seconds.", remaining.Seconds())
	}
	app.MemberSync.Mutex.Lock()
	defer app.MemberSync.Mutex.Unlock()

	memberValueRanges, err := getMemberValueRanges(ctx, app)
	if err != nil {
		return fmt.Errorf("Failed to update members: %v", err)
	}

	formattedMemberStructs := getMemberStructs(memberValueRanges)

	transaction, err := app.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to create a transaction: %v", err)
	}
	for _, each := range formattedMemberStructs {
		err := upsertMember(ctx, each, transaction)
		if err != nil {
			return err
		}
	}
	transaction.Commit()

	app.MemberSync.LastUpdated = time.Now()
	return nil
}

// takes in a formatted member struct and a transaction and upserts their row
// checks if a member with the same name exists
// if they don't, insert them
// otherwise, update
func upsertMember(ctx context.Context, member Member, transaction *sqlx.Tx) error {
	result := Member{}
	err := transaction.GetContext(
		ctx, &result,
		"SELECT * from members WHERE first_name = ? AND last_name = ? LIMIT 1",
		member.Firstname, member.Lastname,
	)
	if err == sql.ErrNoRows {
		_, insertErr := transaction.NamedExecContext(
			ctx,
			`INSERT INTO members
			(first_name, nickname, middle_name, last_name, all_hours, term_hours, grad_year, class, strikes, personal_email, school_email, phone_number, shirt_size, paid_dues)
			VALUES
			(:first_name, :nickname, :middle_name, :last_name, :all_hours, :term_hours, :grad_year, :class, :strikes, :personal_email, :school_email, :phone_number, :shirt_size, :paid_dues)`,
			member,
		)
		if insertErr != nil {
			return fmt.Errorf("Issue inserting member during upsert: %v", insertErr)
		}
	} else if err != nil {
		return fmt.Errorf("Issue upserting member: %v", err)
	} else {
		member.ID = result.ID // to update the correct row based on primary key (id)
		_, updateErr := transaction.NamedExecContext(
			ctx,
			`UPDATE members SET 
			first_name=:first_name, nickname=:nickname, middle_name=:middle_name, last_name=:last_name, all_hours=:all_hours, term_hours=:term_hours, grad_year=:grad_year, class=:class, strikes=:strikes, personal_email=:personal_email, school_email=:school_email, phone_number=:phone_number, shirt_size=:shirt_size, paid_dues=:paid_dues
			WHERE id=:id`,
			member,
		)
		if updateErr != nil {
			return fmt.Errorf("Issue updating member during upsert: %v", updateErr)
		}
	}
	return nil
}

// fetches and returns google sheets api value ranges (unformatted)
func getMemberValueRanges(ctx context.Context, app *internal.App) ([]*sheets.ValueRange, error) {
	r := app.Config.MembersSheetRanges
	data, err := app.GoogleServices.Sheets.Spreadsheets.Values.BatchGet(app.Config.SpreadsheetID).Ranges(
		r.Names,
		r.AllHours,
		r.TermHours,
		r.GradYear,
		r.Class,
		r.Strikes,
		r.PersonalEmail,
		r.SchoolEmail,
		r.PhoneNumber,
		r.ShirtSizes,
		r.PaidDues,
	).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to batch get spreadsheet ranges: %v", err)
	}
	return data.ValueRanges, nil
}

// takes the api call value ranges and turns them into an array of member structs
func getMemberStructs(memberValueRanges []*sheets.ValueRange) []Member {
	// gets length based on the length of the names column
	memberValueRangesLength := len(memberValueRanges[0].Values)
	formattedMemberArray := make([]Member, memberValueRangesLength)

	normalizedNames := genericutils.NormalizeStringValues(memberValueRanges[0].Values, memberValueRangesLength)
	normalizedAllHours := genericutils.NormalizeFloatValues(memberValueRanges[1].Values, memberValueRangesLength)
	normalizedTermHours := genericutils.NormalizeFloatValues(memberValueRanges[2].Values, memberValueRangesLength)
	normalizedGradYears := genericutils.NormalizeIntValues(memberValueRanges[3].Values, memberValueRangesLength)
	normalizedClasses := genericutils.NormalizeStringValues(memberValueRanges[4].Values, memberValueRangesLength)
	normalizedStrikes := genericutils.NormalizeIntValues(memberValueRanges[5].Values, memberValueRangesLength)
	normalizedPersonalEmails := genericutils.NormalizeStringValues(memberValueRanges[6].Values, memberValueRangesLength)
	normalizedSchoolEmails := genericutils.NormalizeStringValues(memberValueRanges[7].Values, memberValueRangesLength)
	normalizedPhoneNumbers := genericutils.NormalizeStringValues(memberValueRanges[8].Values, memberValueRangesLength)
	normalizedShirtSizes := genericutils.NormalizeStringValues(memberValueRanges[9].Values, memberValueRangesLength)
	normalizedPaidDues := genericutils.NormalizeBoolValues(memberValueRanges[10].Values, memberValueRangesLength)

	for i := range memberValueRangesLength - 1 {
		name := NewName(normalizedNames[i])

		formattedMemberArray[i] = Member{
			Firstname:     name.First,
			Nickname:      name.Nick,
			Middlename:    name.Middle,
			Lastname:      name.Last,
			AllHours:      normalizedAllHours[i],
			TermHours:     normalizedTermHours[i],
			GradYear:      normalizedGradYears[i],
			Class:         normalizedClasses[i],
			PersonalEmail: normalizedPersonalEmails[i],
			SchoolEmail:   normalizedSchoolEmails[i],
			PhoneNumber:   normalizedPhoneNumbers[i],
			Strikes:       normalizedStrikes[i],
			ShirtSize:     normalizedShirtSizes[i],
			PaidDues:      normalizedPaidDues[i],
			ID:            -1,
		}
	}

	return formattedMemberArray
}
