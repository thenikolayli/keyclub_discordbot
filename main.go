package main

import (
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/memberutils"
	"log"
)

func main() {
	config.LoadConfig()
	defer config.DB.Close()

	// member, err := hoursutils.GetAllRanks(2027, 5, 0, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	member, err := memberutils.GetMember("badiang willian", 0, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	if err != nil {
		log.Fatalf("Something happened: %v", err)
	}
	log.Printf("Member: %+v", member)
}
