package commands

import (
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/memberutils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var RefreshCommand = &discordgo.ApplicationCommand{
	Name:        "refresh",
	Description: "Refresh the members database from the hours sheet",
}

func RefreshHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := memberutils.SyncMembersFromSheet(config.HoursUpdateTimeout, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	if err != nil {
		// UpdateMembers uses an error to signal the "no-op" case (rate limiting).
		if strings.HasPrefix(err.Error(), "Not enough time has passed") {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Members not updated",
							Description: err.Error(),
							Color:       0xf1c232,
							Footer: &discordgo.MessageEmbedFooter{
								Text: genericutils.GetFormattedLastUpdated(config.HoursLastUpdated),
							},
						},
					},
				},
			})
			return
		}
		genericutils.SendErrorErrorEmbed(
			"Error refreshing members",
			err,
			"Baaaaka!",
			session, interaction,
		)
		return
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Members updated",
					Description: "Successfully refreshed members from the hours sheet.",
					Color:       0xc6eb34,
					Footer: &discordgo.MessageEmbedFooter{
						Text: genericutils.GetFormattedLastUpdated(config.HoursLastUpdated),
					},
				},
			},
		},
	})
}
