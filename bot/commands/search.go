package commands

import (
	"context"
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
		slots := 0
		if len(interaction.ApplicationCommandData().Options) > 0 {
			slots = int(interaction.ApplicationCommandData().Options[0].IntValue())
		}
		events, err := eventutils.SearchEvents(ctx, app, slots)
		if err != nil {
			genericutils.SendErrorErrorEmbed(
				"Error searching events",
				err,
				"Baaaaka!",
				session, interaction,
			)
			return
		}
		if len(events) == 0 {
			genericutils.SendStringErrorEmbed(
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
			description.WriteString(fmt.Sprintf("%v slots open\n", event.NofSlots))
			description.WriteString(fmt.Sprintf("Sign up URL: %s\n\n", event.SignUpUrl))
		}

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       fmt.Sprintf("%v Events Found", len(events)),
						Color:       0xc6eb34,
						Description: description.String(),
						Footer: &discordgo.MessageEmbedFooter{
							Text: genericutils.GetFormattedLastUpdated(app.EventSync.LastUpdated),
						},
					},
				},
			},
		})
	}
}
