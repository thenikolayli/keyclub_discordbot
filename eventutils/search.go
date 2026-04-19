package eventutils

import (
	"context"
	"fmt"
	"keyclubDiscordBot/internal"
	"time"
)

func SearchEvents(ctx context.Context, app *internal.App, slots int) ([]Event, error) {
	SyncEvents(ctx, app)
	today := time.Now().Format(time.DateOnly)
	end := time.Now().AddDate(0, 1, 0).Format(time.DateOnly)

	const query = `
SELECT *
FROM events
WHERE date >= ?
AND date <= ?
AND (n_of_slots - n_of_volunteers) >= ?
ORDER BY date ASC`

	var events []Event
	if err := app.DB.SelectContext(ctx, &events, query, today, end, slots); err != nil {
		return nil, fmt.Errorf("searching events: %w", err)
	}
	return events, nil
}
