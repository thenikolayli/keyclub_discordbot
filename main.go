package main

import (
	"keyclubDiscordBot/bot"
	"keyclubDiscordBot/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	config.LoadConfig()

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
