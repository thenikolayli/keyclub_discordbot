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

	// match first_name to formattedName.Firstname first
	// then nickname to formattedName.Nickname
	// then first_name to formattedName.Lastname in case the user inputted last name first
	err := database.GetContext(
		config.Context, &result,
		`SELECT * FROM members WHERE first_name = ? OR nickname = ? OR first_name = ? LIMIT 1`,
		formattedName.Firstname, formattedName.Nickname, formattedName.Lastname,
	)
	if err == sql.ErrNoRows {
		return Member{}, fmt.Errorf("No member found with the name %v", name)
	}
	if err != nil {
		return Member{}, fmt.Errorf("Error getting member hours: %v", err)
	}

	return result, nil
}
