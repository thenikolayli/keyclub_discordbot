package commands

import (
	"fmt"
	"keyclubDiscordBot/config"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/memberutils"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var TermRanksCommand = &discordgo.ApplicationCommand{
	Name:        "termranks",
	Description: "Check a volunteer's ranks",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "gradyear",
			Description: "Graduation year of the group to view ranks for",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "topn",
			Description: "Number of top volunteers to view ranks for",
			Required:    false,
		},
	},
}

func TermRanksHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	gradYear := interaction.ApplicationCommandData().Options[0].StringValue()
	topNInt := config.DefaultRankTopN                          // default to top 5 ranks
	if len(interaction.ApplicationCommandData().Options) > 1 { // if topN was given
		topNIntContender, err := strconv.Atoi(interaction.ApplicationCommandData().Options[1].StringValue())
		if err == nil {
			topNInt = topNIntContender
		}
	}
	gradYearInt, err := strconv.Atoi(gradYear)
	if err != nil {
		genericutils.SendErrorErrorEmbed(
			"Error parsing graduation year",
			fmt.Errorf("Issue parsing graduation year: %w", err),
			fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("Jan 2 2006 15:04:05")),
			session,
			interaction,
		)
		return
	}

	if topNInt == 0 || topNInt < -1 {
		genericutils.SendStringErrorEmbed(
			"Error fetching ranks",
			"This command does not take most negative topN int values.",
			fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("Jan 2 2006 15:04:05")),
			session,
			interaction,
		)
		return
	}

	members, err := memberutils.GetTermRanks(gradYearInt, topNInt, config.HoursUpdateTimeout, &config.HoursLastUpdated, config.GoogleServices.Sheets, config.DB)
	if err != nil {
		genericutils.SendErrorErrorEmbed(
			"Error fetching ranks",
			fmt.Errorf("Issue fetching ranks: %w", err),
			fmt.Sprintf("Last updated: %v", config.HoursLastUpdated.Format("Jan 2 2006 15:04:05")),
			session,
			interaction,
		)
		return
	}
	indexes := []string{}
	names := []string{}
	hours := []string{}

	for index, member := range members {
		formattedMember := member.Format()
		indexes = append(indexes, fmt.Sprintf("%v.", index+1))
		names = append(names, formattedMember.Name)
		hours = append(hours, strconv.FormatFloat(member.AllHours, 'f', 2, 64))
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Top %v Ranks for Term Hours", topNInt),
					Color: 0xc6eb34,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Rank",
							Value:  strings.Join(indexes, "\n"),
							Inline: true,
						},
						{
							Name:   "Names",
							Value:  strings.Join(names, "\n"),
							Inline: true,
						},
						{
							Name:   "Hours",
							Value:  strings.Join(hours, "\n"),
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
