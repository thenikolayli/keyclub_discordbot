package main

import (
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/hoursutils"
	"log"
)

func main() {
	config.LoadConfig()

	err := hoursutils.UpdateMembers(config.GoogleServices)

	if err != nil {
		log.Fatalf("Something happened: %v", err)
	}
}
