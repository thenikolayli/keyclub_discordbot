package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session        *discordgo.Session
	GuildID        string
	registeredCmds []*discordgo.ApplicationCommand
}

func New(token, guildID string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Bot{Session: session, GuildID: guildID}, nil
}

func (bot *Bot) Start() error {
	bot.Session.AddHandler(Router)
	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", s.State.User.Username)
	})

	if err := bot.Session.Open(); err != nil {
		return err
	}

	for _, command := range Commands {
		registered, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, bot.GuildID, command)
		if err != nil {
			bot.Session.Close()
			return err
		}
		bot.registeredCmds = append(bot.registeredCmds, registered)
	}
	return nil
}

func (bot *Bot) Stop() {
	for _, cmd := range bot.registeredCmds {
		bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, bot.GuildID, cmd.ID)
	}
	bot.Session.Close()
}
