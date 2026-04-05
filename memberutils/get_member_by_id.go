package memberutils

import (
	"database/sql"
	"fmt"
	"time"

	"keyclubDiscordBot/config"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// takes in a member name and returns their member struct
// if the hours haven't been updated in their specified timeout, it will update the hours before returning the member struct
func GetMemberByDiscordId(discordId string, hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) (Member, error) {
	// attempts to update members if enough time has passed since the last update
	UpdateMembers(hoursUpdateTimeout, hoursLastUpdated, sheetsService, database)
	var result Member

	// match first_name to formattedName.Firstname first
	// then nickname to formattedName.Nickname
	// then first_name to formattedName.Lastname in case the user inputted last name first
	err := database.GetContext(
		config.Context, &result,
		`SELECT * FROM members WHERE discord_id = ? LIMIT 1`,
		discordId,
	)
	if err == sql.ErrNoRows {
		return Member{}, fmt.Errorf("No member found with the Discord ID %v", discordId)
	}
	if err != nil {
		return Member{}, fmt.Errorf("Error getting member hours: %v", err)
	}

	return result, nil
}
