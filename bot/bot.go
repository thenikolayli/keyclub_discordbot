package bot

import (
	"keyclubDiscordBot/internal"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	GuildID string
	App     *internal.App
}

func New(app *internal.App) (*Bot, error) {
	session, err := discordgo.New("Bot " + app.Config.DiscordToken)
	if err != nil {
		return nil, err
	}
	return &Bot{Session: session, GuildID: app.Config.GuildID, App: app}, nil
}

func (bot *Bot) Start() error {
	bot.Session.AddHandler(bot.router())
	bot.Session.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Printf("Logged in as %s", session.State.User.Username)
	})

	if err := bot.Session.Open(); err != nil {
		return err
	}

	_, err := bot.Session.ApplicationCommandBulkOverwrite(bot.Session.State.User.ID, bot.GuildID, Commands)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) Stop() {
	bot.Session.Close()
}
