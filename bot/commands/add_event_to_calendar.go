package commands

import (
	"context"
	"fmt"
	"keyclubDiscordBot/eventutils"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/internal"

	"github.com/bwmarrin/discordgo"
)

var AddEventToCalendarCommand = &discordgo.ApplicationCommand{
	Name:        "addeventtocalendar",
	Description: "Add an event to the calendar",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "url",
			Description: "URL of the attendance document",
			Required:    true,
		},
	},
}

func AddEventToCalendarHandler(app *internal.App) func(context.Context, *discordgo.Session, *discordgo.InteractionCreate) {
	return func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		}); err != nil {
			return
		}

		attendanceDocURL := interaction.ApplicationCommandData().Options[0].StringValue()
		attendanceDocID := eventutils.DocsUrlToId(attendanceDocURL)

		event, err := eventutils.AddEventToCalendar(ctx, app, attendanceDocID)
		if err != nil {
			_ = genericutils.EditInteractionErrorEmbed(
				"Error adding event to calendar",
				err,
				"Baaaaka!",
				session, interaction,
			)
			return
		}

		_ = genericutils.EditInteractionEmbeds(session, interaction, []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("%s added to calendar", event.Summary),
				Color:       0xc6eb34,
				Description: fmt.Sprintf("View the event at \n%v", event.HtmlLink),
				Footer: &discordgo.MessageEmbedFooter{
					Text: genericutils.GetFormattedLastUpdated(app.EventSync.LastUpdated),
				},
			},
		})
	}
}
