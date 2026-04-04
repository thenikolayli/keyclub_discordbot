package config

// information about the event sheet ranges
type EventsSheetRangesType struct {
	Names         string
	Nicknames     string
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
type MembersSheetRangesType struct {
	Names         string
	Dates         string
	StartTimes    string
	EndTimes      string
	Addresses     string
	NofSlots      string
	NofVolunteers string
	Tags          string
}

// information about the eventsmembers sheet ranges
type EventsMembersSheetRangesType struct {
	SheetName       string
	Events          string
	Members         string
	MemberNicknames string
}
