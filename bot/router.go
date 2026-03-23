package bot

import "github.com/bwmarrin/discordgo"

func Router(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if h, ok := CommandHandlers[interaction.ApplicationCommandData().Name]; ok {
		h(session, interaction)
	}
}
