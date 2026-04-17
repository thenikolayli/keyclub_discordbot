package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	_ "modernc.org/sqlite"
)

var (
	SpreadsheetID     string
	CalendarID        string
	GoogleAuthKeyPath string

	MembersSheetRanges MembersSheetRangesType = MembersSheetRangesType{
		SheetName:     "2025-2026 Members",
		Names:         "2025-2026 Members!A1:A",
		AllHours:      "2025-2026 Members!B1:B",
		TermHours:     "2025-2026 Members!C1:C",
		GradYear:      "2025-2026 Members!D1:D",
		Class:         "2025-2026 Members!E1:E",
		Strikes:       "2025-2026 Members!F1:F",
		PersonalEmail: "2025-2026 Members!G1:G",
		SchoolEmail:   "2025-2026 Members!H1:H",
		PhoneNumber:   "2025-2026 Members!I1:I",
		ShirtSizes:    "2025-2026 Members!J1:J",
		PaidDues:      "2025-2026 Members!K1:K",
	}

	EventsMembersSheetRanges EventsMembersSheetRangesType = EventsMembersSheetRangesType{
		SheetName: "2025-2026 EventsMembers",
		Events:    "2025-2026 EventsMembers!A1:A",
		Members:   "2025-2026 EventsMembers!A1:ZZ1",
	}

	EventsSheetRanges EventsSheetRangesType = EventsSheetRangesType{
		SheetName:     "2025-2026 Events",
		Events:        "2025-2026 Events!A1:A",
		Dates:         "2025-2026 Events!B1:B",
		StartTimes:    "2025-2026 Events!C1:C",
		EndTimes:      "2025-2026 Events!D1:D",
		Addresses:     "2025-2026 Events!E1:E",
		NofSlots:      "2025-2026 Events!F1:F",
		NofVolunteers: "2025-2026 Events!G1:G",
		TotalHours:    "2025-2026 Events!H1:H",
		Leaders:       "2025-2026 Events!I1:I",
		MadeBy:        "2025-2026 Events!J1:J",
	}

	// EventsSheetRanges

	LeaderRoleId  string
	OfficerRoleId string

	Officers []string

	DB *sqlx.DB

	HoursUpdateTimeout  float64   = 60 * 60
	EventsUpdateTimeout float64   = 60 * 5
	HoursLastUpdated    time.Time = time.Date(2026, time.January, 1, 01, 01, 0, 0, time.UTC)
	EventsLastUpdated   time.Time = time.Date(2026, time.January, 1, 01, 01, 0, 0, time.UTC)

	DefaultRankTopN int = 5

	Context        context.Context
	GoogleServices *GoogleServicesType
	DiscordToken   string
	GuildID        string
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	SpreadsheetID = os.Getenv("HOURS_SPREADSHEET_ID")
	GoogleAuthKeyPath = os.Getenv("GOOGLE_AUTH_KEY_PATH")
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	GuildID = os.Getenv("GUILD_ID")
	CalendarID = os.Getenv("CALENDAR_ID")

	LeaderRoleId = os.Getenv("LEADER_ROLE_ID")
	OfficerRoleId = os.Getenv("OFFICER_ROLE_ID")

	err := prepDatabase()
	if err != nil {
		log.Fatalf("Failed to prepare database: %v", err)
	}

	Context = context.Background()
	GoogleServicesContender, err := getGoogleServices(Context, GoogleAuthKeyPath)
	if err != nil {
		log.Fatalf("Issue getting Google services: %v", err)
	}
	GoogleServices = GoogleServicesContender

	if rawOfficers := os.Getenv("OFFICERS"); rawOfficers != "" {
		Officers = strings.Split(rawOfficers, ", ")
	}
}

// prepares the database and runs migrations
func prepDatabase() error {
	DBContender, err := sqlx.Connect("sqlite", "db.sqlite3")
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}
	DB = DBContender

	migrations, migrationErr := migrate.New(
		"file://migrations",
		"sqlite://db.sqlite3",
	)
	if migrationErr != nil {
		return fmt.Errorf("Failed to initialize migrations: %w", migrationErr)
	}
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}
	return nil
}

// returns a GoogleServices struct, containing the services used to interact with the google apis
func getGoogleServices(ctx context.Context, keyFilePath string) (*GoogleServicesType, error) {
	clientOption, err := getClientOption(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get client option: %v", err)
	}

	docsService, err := docs.NewService(ctx, clientOption)
	if err != nil {
		return nil, fmt.Errorf("google.GetGoogleServices: %w", err)
	}
	sheetsService, err := sheets.NewService(ctx, clientOption)
	if err != nil {
		return nil, fmt.Errorf("google.GetGoogleServices: %w", err)
	}
	calendarService, err := calendar.NewService(ctx, clientOption)
	if err != nil {
		return nil, fmt.Errorf("google.GetGoogleServices: %w", err)
	}

	return &GoogleServicesType{
		Docs:     docsService,
		Sheets:   sheetsService,
		Calendar: calendarService,
	}, nil
}

// uses the google_auth_key.json file to create client options
// this is used to get google services later
func getClientOption(keyFilepath string) (option.ClientOption, error) {
	if _, err := os.Stat(keyFilepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Service account key file not found: %v", err)
	}

	return option.WithAuthCredentialsFile(option.ServiceAccount, keyFilepath), nil
}
