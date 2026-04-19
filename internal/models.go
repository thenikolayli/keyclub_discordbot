package internal

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
