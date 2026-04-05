package commands

import (
	"fmt"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/memberutils"

	"github.com/bwmarrin/discordgo"
)

var HoursCommand = &discordgo.ApplicationCommand{
	Name:        "hours",
	Description: "Check a volunteer's hours",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "Name of the person to view hours for",
			Required:    false, // fallback to user discordid if name not provided
		},
	},
}

func HoursHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	name := interaction.ApplicationCommandData().Options[0].StringValue()
	member, err := memberutils.GetMember(name, config.HoursUpdateTimeout, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	// if member not found, try via user id
	if err != nil {
		member, err = memberutils.GetMemberByDiscordId(interaction.Member.User.ID, config.HoursUpdateTimeout, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
		if err != nil {
			genericutils.SendStringErrorEmbed(
				"Member not found",
				fmt.Sprintf(`Could not find a member with the name "%v" or Discord ID %v.`, name, interaction.Member.User.ID),
				genericutils.GetFormattedLastUpdated(config.HoursLastUpdated),
				session,
				interaction,
			)
			return
		}
	}

	formattedMember := member.Format()
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("%v's Hours", formattedMember.Name),
					Color: 0xc6eb34,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "All time hours",
							Value:  fmt.Sprintf("%v hours", formattedMember.AllHours),
							Inline: true,
						},
						{
							Name:   "Term hours",
							Value:  fmt.Sprintf("%v hours", formattedMember.TermHours),
							Inline: true,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: genericutils.GetFormattedLastUpdated(config.HoursLastUpdated),
					},
				},
			},
		},
	})
}
