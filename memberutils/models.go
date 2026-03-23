package memberutils

import (
	"strings"
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
	ClassYear     int     `db:"class_year"`
	Strikes       int     `db:"strikes"`
	PersonalEmail string  `db:"personal_email"`
	SchoolEmail   string  `db:"school_email"`
	PhoneNumber   string  `db:"phone_number"`
	ShirtSize     string  `db:"shirt_size"`
	PaidDues      bool    `db:"paid_dues"`
	DiscordID     string  `db:"discord_id"`
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

	// if only one name was provided, assume that it's either the nickname or the first name
	if len(nameParts) == 1 {
		return Name{
			Firstname: strings.ToLower(nameParts[0]),
			Lastname:  strings.ToLower(nameParts[0]),
			Nickname:  strings.ToLower(nameParts[0]),
		}
	} else { // otherwise, assume the full name was provided
		return Name{
			Firstname: strings.ToLower(nameParts[0]),
			Lastname:  strings.ToLower(nameParts[1]),
			Nickname:  strings.ToLower(nameParts[0]),
		}
	}
}
