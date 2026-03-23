package commands

import (
	"fmt"
	"keyclubDiscordBot/config"
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
			Required:    true,
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
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Member not found",
							Description: fmt.Sprintf(`Could not find a member with the name "%v" or Discord ID %v.`, name, interaction.Member.User.ID),
							Color:       0xd1244c,
							Footer: &discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("Jan 2 2006 15:04:05")),
							},
						},
					},
				},
			})
			return
		}
	}

	header := ""
	if member.Nickname != "" {
		header = fmt.Sprintf(`Hours for %v "%v" %v`, member.Firstname, member.Nickname, member.Lastname)
	} else {
		header = fmt.Sprintf("Hours for %v %v", member.Firstname, member.Lastname)
	}
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: header,
					Color: 0xc6eb34,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "All time hours",
							Value:  fmt.Sprintf("%v hours", member.AllHours),
							Inline: true,
						},
						{
							Name:   "Term hours",
							Value:  fmt.Sprintf("%v hours", member.TermHours),
							Inline: true,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("2006-01-02 15:04:05")),
					},
				},
			},
		},
	})

}
