package eventutils

// logs the event by adding it to the google calendar
// check if the sign up sheet has hours calculated
// if no, calculate them then move on, if yes, move on
// finds first next empty col in the hours sheet and adds it on there
// returns Event struct and an array of members with hours logged and not logged
func LogEvent(signUpDocUrl string) (LogEventResponse, error) {
	// first, fetch the tables
	// get the metadata table and extract event info from it
	// get the other tables and check if the hours have been calculated
	// if no, calculate them and write them, if yes, move on
	// find next empty column in the EventsMembers sheet and write to it
	// find next empty row in Events and write to it

}
