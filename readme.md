# Key Club Discord Bot v2.0.1

This is t yeah i'll finish this later...

## The plan

Here are some of the features I'd like to add.

## _Member Functions_

- `/me`
  - `userId discord user id` - Given or inferred
  - Takes a userID or infers it from the message and returns info about the member. Info such as term hours, all time hours, grad year, class rank, and events signed up to (this feature might come later).
- `/search` - [ ]
  - `start`, `end`, `slots open`, `leader slots`, `member slots`
  - Takes the values above and finds an event that matches every condition.
- `/refresh` - [ ] MAYBE LATER... low priority.
  - Refreshes the cache
- `/termRanks`
  - `gradYear int`
  - Returns ranks based on term hours
- `/allRanks`
  - `gradYear int`
  - Returns ranks based on all time hours

## _Leader Functions_

- `/member`
  - `name string`
  - Takes a name and returns info about the member (only runs in the private leaders only channel).
  - Returns the `Member` struct
- `/addEvent` - [ ]
  - `signUpUrl string`
  - Takes the sign up google doc url and adds the event to the calendar.

## _Officer Functions_

- `/logEvent` - [ ]
  - `signUpUrl string`
  - Takes the sign up google doc url and logs the event.
- `/getMonthData` - [ ]
  - `month string, year string`
  - Takes the month and year and finds all the events that took place that month.
  - Returns `monthInfo MonthInfo`
  - The `MonthInfo` struct has the fields `Events []Event`, `NofEvents int`, `TotalHours float64`
  - The `Event` struct has the fields `Date (date)`, `StartTime`, `EndTime`, `Address`, `Leaders`, `NofVolunteers`

## Relationships

- Events to Members
  - Many to Many: many members can volunteer at many events
- Events to Leaders
  - Many to Many: many leaders can lead many events
