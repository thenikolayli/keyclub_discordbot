package memberutils

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// struct to represent a member
type Member struct {
	ID            int     `db:"id"`
	Firstname     string  `db:"first_name"`
	Lastname      string  `db:"last_name"`
	Nickname      string  `db:"nickname"`
	AllHours      float64 `db:"all_hours"`
	TermHours     float64 `db:"term_hours"`
	GradYear      int     `db:"grad_year"`
	Class         string  `db:"class"`
	Strikes       int     `db:"strikes"`
	PersonalEmail string  `db:"personal_email"`
	SchoolEmail   string  `db:"school_email"`
	PhoneNumber   string  `db:"phone_number"`
	ShirtSize     string  `db:"shirt_size"`
	PaidDues      bool    `db:"paid_dues"`
	DiscordID     string  `db:"discord_id"`
}

type FormattedMember struct {
	Name          string
	AllHours      float64
	TermHours     float64
	GradYear      int
	Class         string
	Strikes       int
	PersonalEmail string
	SchoolEmail   string
	PhoneNumber   string
	ShirtSize     string
	PaidDues      bool
	DiscordID     string
}

// formats a member for a more readable output, such as for the member lookup command
func (member Member) Format() FormattedMember {
	name := cases.Title(language.English).String(member.Firstname) + " " + cases.Title(language.English).String(member.Lastname)
	if member.Nickname != "" {
		name = cases.Title(language.English).String(member.Firstname) + ` "` + cases.Title(language.English).String(member.Nickname) + `" ` + cases.Title(language.English).String(member.Lastname)
	}
	if member.PersonalEmail == "" {
		member.PersonalEmail = "N/A"
	}
	if member.SchoolEmail == "" {
		member.SchoolEmail = "N/A"
	}
	if member.PhoneNumber == "" {
		member.PhoneNumber = "N/A"
	} else {
		member.PhoneNumber = formatPhoneNumber(member.PhoneNumber)
	}
	if member.ShirtSize == "" {
		member.ShirtSize = "N/A"
	}
	if member.DiscordID == "" {
		member.DiscordID = "N/A"
	}
	return FormattedMember{
		Name:          name,
		AllHours:      member.AllHours,
		TermHours:     member.TermHours,
		GradYear:      member.GradYear,
		Class:         member.Class,
		Strikes:       member.Strikes,
		PersonalEmail: member.PersonalEmail,
		SchoolEmail:   member.SchoolEmail,
		PhoneNumber:   member.PhoneNumber,
		ShirtSize:     member.ShirtSize,
		PaidDues:      member.PaidDues,
		DiscordID:     member.DiscordID,
	}
}

// formats phone numbers into this standard format: (XXX) XXX-XXXX
func formatPhoneNumber(phoneNumber string) string {
	cleanNumber := strings.ReplaceAll(phoneNumber, " ", "")
	cleanNumber = strings.ReplaceAll(cleanNumber, "-", "")
	cleanNumber = strings.ReplaceAll(cleanNumber, "(", "")
	cleanNumber = strings.ReplaceAll(cleanNumber, ")", "")

	if len(cleanNumber) == 10 {
		return fmt.Sprintf("(%s) %s-%s", cleanNumber[0:3], cleanNumber[3:6], cleanNumber[6:10])
	} else {
		return phoneNumber // if the number isn't 10 digits, return it as is
	}
}

// struct to represnt a formatted name
type Name struct {
	Firstname string
	Lastname  string
	Nickname  string
}

// creates a new instance of type Name based on a string input
// this is a standalone function because it's often called on a string input, not a member struct
func NewName(name string) Name {
	nameParts := strings.Split(name, " ")

	// if only one name was provided, assume that it's either the fist name or the nickname
	if len(nameParts) == 1 {
		return Name{
			Firstname: strings.ToLower(nameParts[0]),
			Lastname:  "",
			Nickname:  strings.ToLower(nameParts[0]),
		}
	} else { // otherwise, assume the full name was provided
		return Name{
			Firstname: strings.ToLower(nameParts[0]),
			Lastname:  strings.ToLower(nameParts[1]),
			Nickname:  "",
		}
	}
}
