package main

import (
	"context"
	"keyclubDiscordBot/bot"
	"keyclubDiscordBot/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	var app *internal.App
	var discordBot *bot.Bot
	var err error
	defer func() {
		if discordBot != nil {
			discordBot.Stop()
		}
		if app != nil && app.DB != nil {
			app.DB.Close()
		}
		stop()
	}()

	app, err = internal.NewApp(ctx)
	if err != nil {
		log.Printf("Failed to create app: %v", err)
		return
	}

	discordBot, err = bot.New(app)
	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return
	}
	if err := discordBot.Start(); err != nil {
		log.Printf("Failed to start bot: %v", err)
		discordBot.Stop()
		return
	}

	<-ctx.Done()
}
