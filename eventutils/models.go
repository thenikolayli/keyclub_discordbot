package eventutils

// struct to represent a logged event
type Event struct {
	ID            int     `db:"id"`
	Name          string  `db:"name"`
	Date          string  `db:"date"`
	StartTime     string  `db:"start_time"`
	EndTime       string  `db:"end_time"`
	Address       string  `db:"address"`
	NofSlots      int     `db:"n_of_slots"`
	NofVolunteers int     `db:"n_of_volunteers"`
	TotalHours    float64 `db:"total_hours"`
	Tags          string  `db:"tags"`
}

// struct to represent the intermediary table for many-to-many relationship between events and members
type EventMember struct {
	EventID  int `db:"event_id"`
	MemberID int `db:"member_id"`
}

// essentially the same as EventMember, but for leaders
type EventLeader struct {
	EventID  int `db:"event_id"`
	MemberID int `db:"member_id"`
}

// contains information about the event logged
type LogEventResponse struct {
	Event            Event
	MembersLogged    []MemberAttendance
	MembersNotLogged []MemberAttendance
}

// to represent a row in a sign up doc
// hours start index is the index of the hours cell, it's where you put the calculated hours
type MemberAttendance struct {
	Name            string
	Hours           float64
	HoursStartIndex int
	HoursEndIndex   int
	ColumnFound     bool
}
