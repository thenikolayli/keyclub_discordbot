package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type Config struct {
	SpreadsheetID     string
	CalendarID        string
	GoogleAuthKeyPath string
	DiscordToken      string
	GuildID           string
	OfficerRoleID     string
	LeaderRoleID      string
	DefaultRankTopN   int

	MemberSyncTimeout time.Duration
	EventSyncTimeout  time.Duration
	// DiscordCommandTimeout caps how long a single slash command may run (Sheets/DB/API work).
	DiscordCommandTimeout time.Duration

	MembersSheetRanges       MembersSheetRangesType
	EventsMembersSheetRanges EventsMembersSheetRangesType
	EventsSheetRanges        EventsSheetRangesType
	Officers                 []string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Failed to load .env: %w", err)
	}

	return Config{
		SpreadsheetID:     os.Getenv("HOURS_SPREADSHEET_ID"),
		CalendarID:        os.Getenv("CALENDAR_ID"),
		GoogleAuthKeyPath: os.Getenv("GOOGLE_AUTH_KEY_PATH"),
		DiscordToken:      os.Getenv("DISCORD_TOKEN"),
		GuildID:           os.Getenv("GUILD_ID"),
		OfficerRoleID:     os.Getenv("OFFICER_ROLE_ID"),
		LeaderRoleID:      os.Getenv("LEADER_ROLE_ID"),
		DefaultRankTopN:   5,

		MemberSyncTimeout:     time.Minute * 60,
		EventSyncTimeout:      time.Minute * 10,
		DiscordCommandTimeout: time.Second * 30,

		MembersSheetRanges: MembersSheetRangesType{
			SheetName:     "2025-2026 Members",
			Names:         "2025-2026 Members!A2:A",
			AllHours:      "2025-2026 Members!B2:B",
			TermHours:     "2025-2026 Members!C2:C",
			GradYear:      "2025-2026 Members!D2:D",
			Class:         "2025-2026 Members!E2:E",
			Strikes:       "2025-2026 Members!F2:F",
			PersonalEmail: "2025-2026 Members!G2:G",
			SchoolEmail:   "2025-2026 Members!H2:H",
			PhoneNumber:   "2025-2026 Members!I2:I",
			ShirtSizes:    "2025-2026 Members!J2:J",
			PaidDues:      "2025-2026 Members!K2:K",
		},
		EventsMembersSheetRanges: EventsMembersSheetRangesType{
			SheetName: "2025-2026 EventsMembers",
			Events:    "2025-2026 EventsMembers!A2:A",
			Members:   "2025-2026 EventsMembers!B1:ZZ1",
		},
		EventsSheetRanges: EventsSheetRangesType{
			SheetName:     "2025-2026 Events",
			Events:        "2025-2026 Events!A2:A",
			Dates:         "2025-2026 Events!B2:B",
			StartTimes:    "2025-2026 Events!C2:C",
			EndTimes:      "2025-2026 Events!D2:D",
			Addresses:     "2025-2026 Events!E2:E",
			NofSlots:      "2025-2026 Events!F2:F",
			NofVolunteers: "2025-2026 Events!G2:G",
			TotalHours:    "2025-2026 Events!H2:H",
			Leaders:       "2025-2026 Events!I2:I",
			MadeBy:        "2025-2026 Events!J2:J",
		},
		Officers: strings.Split(os.Getenv("OFFICERS"), ", "),
	}, nil
}
