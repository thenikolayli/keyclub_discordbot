package genericutils

import (
	"github.com/bwmarrin/discordgo"
)

// reusable function for sending errors
func SendErrorErrorEmbed(title string, description error, footerText string, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: description.Error(),
					Color:       0xd1244c,
					Footer: &discordgo.MessageEmbedFooter{
						Text: footerText,
					},
				},
			},
		},
	})
}

func SendStringErrorEmbed(title string, description string, footerText string, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: description,
					Color:       0xd1244c,
					Footer: &discordgo.MessageEmbedFooter{
						Text: footerText,
					},
				},
			},
		},
	})
}
