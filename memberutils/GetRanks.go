package memberutils

import (
	"fmt"
	"keyclubDiscordBot/config"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/api/sheets/v4"
)

// returns ranks based on hours for a given grad year, sorted from highest to lowest
func GetAllRanks(gradYear int, topN int, hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) ([]Member, error) {
	// attempts to update members if enough time has passed since the last update
	updateErr := UpdateMembers(hoursUpdateTimeout, hoursLastUpdated, sheetsService, database)
	if updateErr != nil {
		return []Member{}, fmt.Errorf("Failed to update members: %v", updateErr)
	}

	ranks := []Member{}
	err := database.SelectContext(
		config.Context, &ranks,
		"select * from members where grad_year = ? order by all_hours desc",
		gradYear,
	)
	if err != nil {
		return []Member{}, fmt.Errorf("Failed to get ranks: %v", err)
	}
	// topNRanks := []Member{}
	// currentIndex := 0
	// for len(topNRanks) < topN {
	// 	if ranks[currentIndex].Firstname
	// }

	return ranks, nil
}

// returns ranks based on hours for a given grad year, sorted from highest to lowest
func GetTermRanks(gradYear int, topN int, hoursUpdateTimeout float64, hoursLastUpdated *time.Time, sheetsService *sheets.Service, database *sqlx.DB) ([]Member, error) {
	// attempts to update members if enough time has passed since the last update
	updateErr := UpdateMembers(hoursUpdateTimeout, hoursLastUpdated, sheetsService, database)
	if updateErr != nil {
		return []Member{}, fmt.Errorf("Failed to update members: %v", updateErr)
	}

	ranks := []Member{}
	err := database.SelectContext(
		config.Context, &ranks,
		"select * from members where grad_year = ? order by term_hours desc",
		gradYear,
	)
	if err != nil {
		return []Member{}, fmt.Errorf("Failed to get ranks: %v", err)
	}
	return ranks, nil
}
