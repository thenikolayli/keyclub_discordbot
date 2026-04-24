package commands

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"keyclubDiscordBot/internal"
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
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "gradyear",
			Description: "Graduation year of the group to view ranks for",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "topn",
			Description: "Number of top volunteers to view ranks for",
			Required:    false,
		},
	},
}

func TermRanksHandler(app *internal.App) func(context.Context, *discordgo.Session, *discordgo.InteractionCreate) {
	return func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		}); err != nil {
			return
		}

		gradYear := int(interaction.ApplicationCommandData().Options[0].IntValue())
		topN := app.Config.DefaultRankTopN                          // default to top 5 ranks
		if len(interaction.ApplicationCommandData().Options) == 2 { // if topN was given
			topNGiven := int(interaction.ApplicationCommandData().Options[1].IntValue())
			if topNGiven == 0 || topNGiven < -1 {
				_ = genericutils.EditInteractionStringErrorEmbed(
					"Error fetching ranks",
					"This command does not take most negative topN int values.",
					genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
					session,
					interaction,
				)
				return
			}
			topN = topNGiven
		}

		members, err := memberutils.GetTermRanks(ctx, app, gradYear, topN)
		if err != nil {
			_ = genericutils.EditInteractionErrorEmbed(
				"Error fetching ranks",
				fmt.Errorf("Issue fetching ranks: %w", err),
				genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
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
			hours = append(hours, strconv.FormatFloat(member.TermHours, 'f', 2, 64))
		}

		_ = genericutils.EditInteractionEmbeds(session, interaction, []*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("Top %v Ranks for Term Hours", topN),
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
					Text: genericutils.GetFormattedLastUpdated(app.MemberSync.LastUpdated),
				},
			},
		})
	}
}
