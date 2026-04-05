package commands

import (
	"fmt"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/eventutils"
	"keyclubDiscordBot/genericutils"

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

func AddEventToCalendarHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	attendanceDocURL := interaction.ApplicationCommandData().Options[0].StringValue()
	attendanceDocID := eventutils.DocsUrlToId(attendanceDocURL)

	event, err := eventutils.AddEventToCalendar(attendanceDocID)
	if err != nil {
		genericutils.SendErrorErrorEmbed(
			"Error adding event to calendar",
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
					Title:       fmt.Sprintf("%s added to calendar", event.Summary),
					Color:       0xc6eb34,
					Description: fmt.Sprintf("View the event at \n%v", event.HtmlLink),
					Footer: &discordgo.MessageEmbedFooter{
						Text: genericutils.GetFormattedLastUpdated(config.HoursLastUpdated),
					},
				},
			},
		},
	})
}
