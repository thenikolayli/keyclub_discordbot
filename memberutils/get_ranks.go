package memberutils

import (
	"context"
	"fmt"
	"keyclubDiscordBot/internal"
	"slices"
)

// returns ranks based on hours for a given grad year, sorted from highest to lowest
func GetAllRanks(ctx context.Context, app *internal.App, gradYear int, topN int) ([]Member, error) {
	SyncMembersFromSheet(ctx, app)
	ranks := []Member{}
	err := app.DB.SelectContext(
		ctx, &ranks,
		"SELECT * FROM members WHERE grad_year = ? ORDER BY all_hours DESC",
		gradYear,
	)
	if err != nil {
		return []Member{}, fmt.Errorf("Failed to get ranks: %v", err)
	}

	// removes officers
	topNRanks := []Member{}
	if topN != -1 && topN <= len(ranks) {
		currentIndex := 0
		for len(topNRanks) < topN {
			if !slices.Contains(app.Config.Officers, fmt.Sprintf("%s %s", ranks[currentIndex].Firstname, ranks[currentIndex].Lastname)) {
				topNRanks = append(topNRanks, ranks[currentIndex])
			}
			currentIndex++
		}
	} else {
		for _, rank := range ranks {
			if !slices.Contains(app.Config.Officers, fmt.Sprintf("%s %s", rank.Firstname, rank.Lastname)) {
				topNRanks = append(topNRanks, rank)
			}
		}
	}

	return topNRanks, nil
}

// returns ranks based on hours for a given grad year, sorted from highest to lowest
func GetTermRanks(ctx context.Context, app *internal.App, gradYear int, topN int) ([]Member, error) {
	SyncMembersFromSheet(ctx, app)
	ranks := []Member{}
	err := app.DB.SelectContext(
		ctx, &ranks,
		"SELECT * FROM members WHERE grad_year = ? ORDER BY term_hours DESC",
		gradYear,
	)
	if err != nil {
		return []Member{}, fmt.Errorf("Failed to get ranks: %v", err)
	}
	// removes officers
	topNRanks := []Member{}
	if topN != -1 && topN <= len(ranks) {
		currentIndex := 0
		for len(topNRanks) < topN {
			if !slices.Contains(app.Config.Officers, fmt.Sprintf("%s %s", ranks[currentIndex].Firstname, ranks[currentIndex].Lastname)) {
				topNRanks = append(topNRanks, ranks[currentIndex])
			}
			currentIndex++
		}
	} else {
		for _, rank := range ranks {
			if !slices.Contains(app.Config.Officers, fmt.Sprintf("%s %s", rank.Firstname, rank.Lastname)) {
				topNRanks = append(topNRanks, rank)
			}
		}
	}

	return topNRanks, nil
}
