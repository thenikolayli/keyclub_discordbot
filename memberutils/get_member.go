package memberutils

import (
	"context"
	"database/sql"
	"fmt"
	"keyclubDiscordBot/internal"
)

// takes in a member name and returns their member struct
// if the hours haven't been updated in their specified timeout, it will update the hours before returning the member struct
func GetMember(ctx context.Context, app *internal.App, name string) (Member, error) {
	SyncMembersFromSheet(ctx, app)
	formattedName := NewName(name)
	result := Member{}
	var err error

	// if first and last was given, try both, then reverse (if they put in last first)
	// then try by nickname if given, then by first
	if formattedName.First != "" && formattedName.Last != "" {
		err = app.DB.GetContext(
			ctx, &result,
			`SELECT * FROM members WHERE first_name = ? AND last_name = ? LIMIT 1`,
			formattedName.First, formattedName.Last,
		)
		if err == sql.ErrNoRows {
			err = app.DB.GetContext(
				ctx, &result,
				`SELECT * FROM members WHERE first_name = ? AND last_name = ? LIMIT 1`,
				formattedName.Last, formattedName.First,
			)
		}
	} else if formattedName.Nick != "" {
		err = app.DB.GetContext(
			ctx, &result,
			`SELECT * FROM members WHERE nickname = ? OR first_name = ? LIMIT 1`,
			formattedName.Nick, formattedName.First,
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

func GetMemberByDiscordId(ctx context.Context, app *internal.App, discordId string) (Member, error) {
	SyncMembersFromSheet(ctx, app)
	var result Member

	err := app.DB.GetContext(
		ctx, &result,
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
