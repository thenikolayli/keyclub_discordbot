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
		if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		}); err != nil {
			return
		}

		err := memberutils.SyncMembersFromSheet(ctx, app)
		if err != nil {
			// UpdateMembers uses an error to signal the "no-op" case (rate limiting).
			if strings.HasPrefix(err.Error(), "Not enough time has passed") {
				_ = genericutils.EditInteractionErrorEmbed(
					"Members not updated",
					err,
					genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
					session,
					interaction,
				)
			} else {
				_ = genericutils.EditInteractionErrorEmbed(
					"Error refreshing members",
					err,
					"Baaaaka!",
					session, interaction,
				)
			}
			return
		}

		_ = genericutils.EditInteractionEmbeds(session, interaction, []*discordgo.MessageEmbed{
			{
				Title:       "Members updated",
				Description: "Successfully refreshed members from the hours sheet.",
				Color:       0xc6eb34,
				Footer: &discordgo.MessageEmbedFooter{
					Text: genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
				},
			},
		})
	}
}
