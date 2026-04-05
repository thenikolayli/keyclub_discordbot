package memberutils

import (
	"fmt"
	"keyclubDiscordBot/config"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// returns ranks based on hours for a given grad year, sorted from highest to lowest
func GetAllRanks(gradYear int, topN int, hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) ([]Member, error) {
	// attempts to update members if enough time has passed since the last update
	UpdateMembers(hoursUpdateTimeout, hoursLastUpdated, sheetsService, database)

	ranks := []Member{}
	err := database.SelectContext(
		config.Context, &ranks,
		"SELECT * FROM members WHERE grad_year = ? ORDER BY all_hours DESC",
		gradYear,
	)
	if err != nil {
		return []Member{}, fmt.Errorf("Failed to get ranks: %v", err)
	}

	// removes officers
	topNRanks := []Member{}
	if topN != -1 {
		currentIndex := 0
		for len(topNRanks) < topN {
			if !slices.Contains(config.Officers, fmt.Sprintf("%s %s", ranks[currentIndex].Firstname, ranks[currentIndex].Lastname)) {
				topNRanks = append(topNRanks, ranks[currentIndex])
			}
			currentIndex++
		}
	} else {
		for _, rank := range ranks {
			if !slices.Contains(config.Officers, fmt.Sprintf("%s %s", rank.Firstname, rank.Lastname)) {
				topNRanks = append(topNRanks, rank)
			}
		}
	}

	return topNRanks, nil
}

// returns ranks based on hours for a given grad year, sorted from highest to lowest
func GetTermRanks(gradYear int, topN int, hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) ([]Member, error) {
	// attempts to update members if enough time has passed since the last update
	UpdateMembers(hoursUpdateTimeout, hoursLastUpdated, sheetsService, database)

	ranks := []Member{}
	err := database.SelectContext(
		config.Context, &ranks,
		"SELECT * FROM members WHERE grad_year = ? ORDER BY term_hours DESC",
		gradYear,
	)
	if err != nil {
		return []Member{}, fmt.Errorf("Failed to get ranks: %v", err)
	}
	// removes officers
	topNRanks := []Member{}
	if topN != -1 {
		currentIndex := 0
		for len(topNRanks) < topN {
			if !slices.Contains(config.Officers, fmt.Sprintf("%s %s", ranks[currentIndex].Firstname, ranks[currentIndex].Lastname)) {
				topNRanks = append(topNRanks, ranks[currentIndex])
			}
			currentIndex++
		}
	} else {
		for _, rank := range ranks {
			if !slices.Contains(config.Officers, fmt.Sprintf("%s %s", rank.Firstname, rank.Lastname)) {
				topNRanks = append(topNRanks, rank)
			}
		}
	}

	return topNRanks, nil
}
