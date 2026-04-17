package config

import (
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/sheets/v4"
)

// information about the event sheet ranges
type MembersSheetRangesType struct {
	SheetName     string
	Names         string
	AllHours      string
	TermHours     string
	GradYear      string
	Class         string
	Strikes       string
	PersonalEmail string
	SchoolEmail   string
	PhoneNumber   string
	ShirtSizes    string
	PaidDues      string
}

// information about the members sheet ranges
// currently not used
type EventsSheetRangesType struct {
	SheetName     string
	Events        string
	Dates         string
	StartTimes    string
	EndTimes      string
	Addresses     string
	NofSlots      string
	NofVolunteers string
	TotalHours    string
	Leaders       string
	MadeBy        string
}

// information about the eventsmembers sheet ranges
type EventsMembersSheetRangesType struct {
	SheetName string
	Events    string
	Members   string
}

// struct representing an object containing all Google Services
// to be passed between functions that interact with the google APIs
type GoogleServicesType struct {
	Docs     *docs.Service
	Sheets   *sheets.Service
	Calendar *calendar.Service
}
