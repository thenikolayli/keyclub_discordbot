package genericutils

import (
	"github.com/bwmarrin/discordgo"
)

// EditInteractionEmbeds updates the message for a deferred interaction response
// (after InteractionResponseDeferredChannelMessageWithSource).
func EditInteractionEmbeds(session *discordgo.Session, interaction *discordgo.InteractionCreate, embeds []*discordgo.MessageEmbed) error {
	_, err := session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
	return err
}

func EditInteractionStringErrorEmbed(title, description, footerText string, session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return EditInteractionEmbeds(session, interaction, []*discordgo.MessageEmbed{
		{
			Title:       title,
			Description: description,
			Color:       0xd1244c,
			Footer: &discordgo.MessageEmbedFooter{
				Text: footerText,
			},
		},
	})
}

func EditInteractionErrorEmbed(title string, description error, footerText string, session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return EditInteractionStringErrorEmbed(title, description.Error(), footerText, session, interaction)
}
