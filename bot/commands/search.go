package commands

import (
	"context"
	"errors"
	"fmt"
	"keyclubDiscordBot/eventutils"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/internal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var SearchCommand = &discordgo.ApplicationCommand{
	Name:        "search",
	Description: "Finds upcoming events with >= n slots open",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "slots",
			Description: "N of slots open",
			Required:    false,
		},
	},
}

func SearchHandler(app *internal.App) func(context.Context, *discordgo.Session, *discordgo.InteractionCreate) {
	return func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		}); err != nil {
			return
		}

		slots := 1
		if len(interaction.ApplicationCommandData().Options) > 0 {
			slots = int(interaction.ApplicationCommandData().Options[0].IntValue())
		}
		events, err := eventutils.SearchEvents(ctx, app, slots)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				_ = genericutils.EditInteractionStringErrorEmbed(
					"Search timed out",
					"Syncing or searching events took too long. Try again in a moment.",
					genericutils.GetFormattedLastUpdated(app.EventSync.LastUpdated),
					session, interaction,
				)
				return
			}
			_ = genericutils.EditInteractionErrorEmbed(
				"Error searching events",
				err,
				"Baaaaka!",
				session, interaction,
			)
			return
		}
		if len(events) == 0 {
			_ = genericutils.EditInteractionStringErrorEmbed(
				"No events found",
				fmt.Sprintf("No events found with >= %v slots open", slots),
				genericutils.GetFormattedLastUpdated(app.EventSync.LastUpdated),
				session, interaction,
			)
			return
		}

		var description strings.Builder
		for _, event := range events {
			description.WriteString(fmt.Sprintf("**%s**\n", event.Name))
			description.WriteString(fmt.Sprintf("From %s to %s\n", event.StartTime, event.EndTime))
			description.WriteString(fmt.Sprintf("At %s\n", event.Address))
			description.WriteString(fmt.Sprintf("%v slots open\n", event.NofSlots-event.NofVolunteers))
			description.WriteString(fmt.Sprintf("Sign up URL: %s\n\n", event.SignUpUrl))
		}

		_ = genericutils.EditInteractionEmbeds(session, interaction, []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("%v Events Found", len(events)),
				Color:       0xc6eb34,
				Description: description.String(),
				Footer: &discordgo.MessageEmbedFooter{
					Text: genericutils.GetFormattedLastUpdated(app.EventSync.LastUpdated),
				},
			},
		})
	}
}
