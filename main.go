package main

import (
	"keyclubDiscordBot/bot"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/memberutils"
	"log"
	"os"
	"os/signal"
)

func main() {
	config.LoadConfig()

	// member, err := hoursutils.GetAllRanks(2027, 5, 0, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	// member, err := memberutils.GetMember("badiang willian", 0, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	// member, err := memberutils.GetAllRanks(2027, 5, 0, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	// if err != nil {
	// 	log.Fatalf("Something happened: %v", err)
	// }
	// log.Printf("Member: %+v", member)
	memberutils.UpdateMembers(config.HoursUpdateTimeout, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)

	bot, err := bot.New(config.DiscordToken, config.GuildID)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	defer func() {
		config.DB.Close()
		bot.Stop()
	}()

	// waits until an interrupt signal is received to gracefully shut down the bot and close the database connection
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
