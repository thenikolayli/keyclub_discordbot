package bot

import (
	"context"
	"keyclubDiscordBot/bot/commands"
	"keyclubDiscordBot/internal"
	"slices"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler runs one slash command. ctx is cancelled on shutdown or DiscordCommandTimeout.
type CommandHandler func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate)

var Commands = []*discordgo.ApplicationCommand{
	commands.HoursCommand,
	commands.MemberLookupCommand,
	commands.AllRanksCommand,
	commands.TermRanksCommand,
	commands.LogEventCommand,
	commands.AddEventToCalendarCommand,
	commands.RefreshCommand,
	commands.SearchCommand,
}

// passes a function to get role Ids so they are updated when the function is called so they're not empty upon package initialization
func BuildCommandHandlers(app *internal.App) map[string]CommandHandler {
	return map[string]CommandHandler{
		"hours":              commands.HoursHandler(app),
		"member":             requireRole(func() []string { return []string{app.Config.OfficerRoleID, app.Config.LeaderRoleID} }, commands.MemberLookupHandler(app)),
		"allranks":           commands.AllRanksHandler(app),
		"termranks":          commands.TermRanksHandler(app),
		"logevent":           requireRole(func() []string { return []string{app.Config.OfficerRoleID} }, commands.LogEventHandler(app)),
		"addeventtocalendar": requireRole(func() []string { return []string{app.Config.OfficerRoleID, app.Config.LeaderRoleID} }, commands.AddEventToCalendarHandler(app)),
		"refresh":            commands.RefreshHandler(app),
		"search":             commands.SearchHandler(app),
	}
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
func requireRole(getRequiredRoles func() []string, next CommandHandler) CommandHandler {
	return func(ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
		next(ctx, session, interaction)
	}
}
