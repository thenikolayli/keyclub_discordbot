package bot

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) router() func(*discordgo.Session, *discordgo.InteractionCreate) {
	handlers := BuildCommandHandlers(bot.App)
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Type != discordgo.InteractionApplicationCommand {
			return
		}
		if handler, ok := handlers[interaction.ApplicationCommandData().Name]; ok {
			requestCtx, cancel := context.WithTimeout(bot.App.ShutdownCtx, bot.App.Config.DiscordCommandTimeout)
			defer cancel()
			handler(requestCtx, session, interaction)
		}
	}
}
