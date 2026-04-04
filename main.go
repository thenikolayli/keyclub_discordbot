package main

import (
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/eventutils"
)

func main() {
	config.LoadConfig()

	// bot, err := bot.New(config.DiscordToken, config.GuildID)
	// if err != nil {
	// 	log.Fatalf("Failed to create bot: %v", err)
	// }
	// if err := bot.Start(); err != nil {
	// 	log.Fatalf("Failed to start bot: %v", err)
	// }

	// defer func() {
	// 	config.DB.Close()
	// 	bot.Stop()
	// }()

	// // waits until an interrupt signal is received to gracefully shut down the bot and close the database connection
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt)
	// <-stop

	eventutils.LogEvent("https://docs.google.com/document/d/1MlPWssjw_PRoUmASr60jGLA8DNiS6rPbM5TDW_8m1GU/edit?tab=t.0")
}
