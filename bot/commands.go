package bot

import (
	"keyclubDiscordBot/bot/commands"
	"keyclubDiscordBot/config"
	"slices"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	commands.HoursCommand,
}

var CommandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"hours":  commands.HoursHandler,
	"member": requireRole([]string{config.OfficerRoleId, config.LeaderRoleId}, commands.MemberLookupHandler),
}

// checks if a member has a certain role
// loops through every role a member has and checks if they have the role with the given id
func hasRole(member *discordgo.Member, roleIds []string) bool {
	for _, role := range member.Roles {
		if slices.Contains(roleIds, role) {
			return true
		}
	}
	return false
}

// requires a user to have a certain role otherwise it will respond with an error message and not execute the command
// wrapper function that returns a command handler if the user has the required role, otherwise responds with an error message
func requireRole(roleIds []string, next func(*discordgo.Session, *discordgo.InteractionCreate)) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Member == nil || !hasRole(interaction.Member, roleIds) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Unauthorized",
							Description: "You are not authorized to use this command.",
							Color:       0xd1244c,
							Footer: &discordgo.MessageEmbedFooter{
								Text: "xx ;)",
							},
						},
					},
				},
			})
			return
		}
		next(session, interaction)
	}
}
