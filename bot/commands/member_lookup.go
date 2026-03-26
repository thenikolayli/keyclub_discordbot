package commands

import (
	"fmt"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/memberutils"

	"github.com/bwmarrin/discordgo"
)

var MemberLookupCommand = &discordgo.ApplicationCommand{
	Name:        "member",
	Description: "Get info about a member",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "Name of the person to view hours for",
			Required:    true,
		},
	},
}

func MemberLookupHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	name := interaction.ApplicationCommandData().Options[0].StringValue()
	member, err := memberutils.GetMember(name, config.HoursUpdateTimeout, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	// if member not found, respond with an error messasge
	if err != nil {
		genericutils.SendStringErrorEmbed(
			"Member not found",
			fmt.Sprintf(`Could not find a member with the name "%v".`, name),
			fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("Jan 2 2006 15:04:05")),
			session,
			interaction,
		)
		return
	}

	formattedMember := member.Format()

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("%v's Member Info", formattedMember.Name),
					Color: 0xc6eb34,
					Description: fmt.Sprintf(`
					All time hours: *%v*
					Term hours: *%v*
					%v, class of %v
					Phone Number: *%v*
					Personal Email: *%v*
					School Email: *%v*
					Shirt Size: *%v*
					Paid Dues: *%v*
					Strikes: *%v*
					`, formattedMember.AllHours, formattedMember.TermHours, formattedMember.GradYear, formattedMember.Class, formattedMember.PhoneNumber, formattedMember.PersonalEmail, formattedMember.SchoolEmail, formattedMember.ShirtSize, formattedMember.PaidDues, formattedMember.Strikes),
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("Jan 2 2006 15:04:05")),
					},
				},
			},
		},
	})

}
