package commands

import (
	"context"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/internal"
	"keyclubDiscordBot/memberutils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var RefreshCommand = &discordgo.ApplicationCommand{
	Name:        "refresh",
	Description: "Refresh the members database from the hours sheet",
}

func RefreshHandler(app *internal.App) func(context.Context, *discordgo.Session, *discordgo.InteractionCreate) {
	return func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		err := memberutils.SyncMembersFromSheet(ctx, app)
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
									Text: genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
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
							Text: genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
						},
					},
				},
			},
		})
	}
}
