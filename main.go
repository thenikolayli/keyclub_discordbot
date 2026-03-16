package main

import (
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/hoursutils"
	"log"
)

func main() {
	config.LoadConfig()
	defer config.DB.Close()

	err := hoursutils.UpdateMembers(config.GoogleServices, config.DB, config.Context)

	if err != nil {
		log.Fatalf("Something happened: %v", err)
	}
}
