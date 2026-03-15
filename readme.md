# Key Club Discord Bot

This is t yeah i'll finish this later...


## The plan

Here are some of the features I'd like to add.

## _Public Functions_
* `/me`
    * `userId` - Given or inferred
    * Takes in a userID or infers it from the message and returns info about the member. Info such as term hours, all time hours, grad year, class rank, and events signed up to (this feature might come later).
* `/search`
    * `start`, `end`, `slots open`, `leader slots`, `member slots`
    * Takes the values above and finds an event that matches every condition.
* `/refresh`
    * Refreshes the cache

## _Private Functions_
* `/logEvent`
    * `signUpUrl`
    * Takes in the sign up google doc url and logs the event.
* `/findMember`
    * `name`
    * Finds a member and returns valuable info such as the info from `/me` as well as the number of strikes and shirt size.

## Architecture

I will be using Go with DiscordGo.

Database models will be `member`, `event metadata`, and `event signups`. The TTL before refetching upon request for each in order is 60 minutes, 15 minutes, and 15 minutes.

## Utility Functions

- `getTables` - gets tables from a Google Doc
- `getMetadata` - gets metadata from a table
- `getSignups` - gets signups from a list of tables
- `isSignedOut` - checks if event volunteers have been signed out
