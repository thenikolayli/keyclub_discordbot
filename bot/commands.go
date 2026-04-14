package bot

import (
	"keyclubDiscordBot/bot/commands"
	"keyclubDiscordBot/config"
	"slices"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	commands.HoursCommand,
	commands.MemberLookupCommand,
	commands.AllRanksCommand,
	commands.TermRanksCommand,
	commands.LogEventCommand,
	commands.AddEventToCalendarCommand,
	commands.RefreshCommand,
}

// passes a function to get role Ids so they are updated when the function is called so they're not empty upon package initialization
var CommandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"hours":              commands.HoursHandler,
	"member":             requireRole(func() []string { return []string{config.OfficerRoleId, config.LeaderRoleId} }, commands.MemberLookupHandler),
	"allranks":           commands.AllRanksHandler,
	"termranks":          commands.TermRanksHandler,
	"logevent":           requireRole(func() []string { return []string{config.OfficerRoleId} }, commands.LogEventHandler),
	"addeventtocalendar": requireRole(func() []string { return []string{config.OfficerRoleId, config.LeaderRoleId} }, commands.AddEventToCalendarHandler),
	"refresh":            commands.RefreshHandler,
}

// checks if a member has a certain role
// loops through every role a member has and checks if they have the role with the given id
func hasRole(member *discordgo.Member, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		if slices.Contains(member.Roles, requiredRole) {
			return true
		}
	}
	return false
}

// requires a user to have a certain role otherwise it will respond with an error message and not execute the command
// wrapper function that returns a command handler if the user has the required role, otherwise responds with an error message
func requireRole(getRequiredRoles func() []string, next func(*discordgo.Session, *discordgo.InteractionCreate)) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		requiredRoles := getRequiredRoles()
		if interaction.Member == nil || !hasRole(interaction.Member, requiredRoles) {
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
