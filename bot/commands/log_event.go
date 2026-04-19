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

var LogEventCommand = &discordgo.ApplicationCommand{
	Name:        "logevent",
	Description: "Log a volunteer event",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "url",
			Description: "URL of the attendance document",
			Required:    true,
		},
	},
}

func LogEventHandler(app *internal.App) func(context.Context, *discordgo.Session, *discordgo.InteractionCreate) {
	return func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		attendanceDocURL := interaction.ApplicationCommandData().Options[0].StringValue()
		attendanceDocID := eventutils.DocsUrlToId(attendanceDocURL)

		logEventResponse, err := eventutils.LogEvent(ctx, app, attendanceDocID)
		if err != nil {
			genericutils.SendErrorErrorEmbed(
				"Error logging event",
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
						Title: fmt.Sprintf("%s logged", logEventResponse.Event.Name),
						Color: 0xc6eb34,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Members Logged",
								Value:  strings.Join(eventutils.FormatMemberAttendances(logEventResponse.MembersLogged), "\n"),
								Inline: true,
							},
							{
								Name:   "Members Not Logged",
								Value:  strings.Join(eventutils.FormatMemberAttendances(logEventResponse.MembersNotLogged), "\n"),
								Inline: true,
							},
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: genericutils.GetFormattedLastUpdated(app.EventSync.LastUpdated),
						},
					},
				},
			},
		})
	}
}
