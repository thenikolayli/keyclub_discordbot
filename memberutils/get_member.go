package memberutils

import (
	"database/sql"
	"fmt"
	"keyclubDiscordBot/config"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// takes in a member name and returns their member struct
// if the hours haven't been updated in their specified timeout, it will update the hours before returning the member struct
func GetMember(name string, hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) (Member, error) {
	// attempts to update members if enough time has passed since the last update
	UpdateMembers(hoursUpdateTimeout, hoursLastUpdated, sheetsService, database)

	formattedName := NewName(name)
	result := Member{}
	var err error

	// if first and last was given, try both, then reverse (if they put in last first)
	// then try by nickname if given, then by first
	if formattedName.Firstname != "" && formattedName.Lastname != "" {
		err = database.GetContext(
			config.Context, &result,
			`SELECT * FROM members WHERE first_name = ? AND last_name = ? LIMIT 1`,
			formattedName.Firstname, formattedName.Lastname,
		)
		if err == sql.ErrNoRows {
			err = database.GetContext(
				config.Context, &result,
				`SELECT * FROM members WHERE first_name = ? AND last_name = ? LIMIT 1`,
				formattedName.Lastname, formattedName.Firstname,
			)
		}
	} else {
		err = database.GetContext(
			config.Context, &result,
			`SELECT * FROM members WHERE nickname = ? OR first_name = ? LIMIT 1`,
			formattedName.Nickname, formattedName.Firstname,
		)
	}
	if err == sql.ErrNoRows {
		return Member{}, fmt.Errorf("No member found with the name %v", name)
	}
	if err != nil {
		return Member{}, fmt.Errorf("Error getting member hours: %v", err)
	}

	return result, nil
}
