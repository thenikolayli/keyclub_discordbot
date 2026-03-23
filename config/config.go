package config

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var (
	SpreadsheetID     string
	GoogleAuthKeyPath string

	NamesRange         string = "2025-2026 Member Info!A2:A"
	NicknamesRange     string = "2025-2026 Member Info!B2:B"
	AllHoursRange      string = "2025-2026 Member Info!C2:C"
	TermHoursRange     string = "2025-2026 Member Info!D2:D"
	GradYearRange      string = "2025-2026 Member Info!E2:E"
	StrikesRange       string = "2025-2026 Member Info!G2:G"
	ClassYearRange     string = "2025-2026 Member Info!F2:F"
	PersonalEmailRange string = "2025-2026 Member Info!H2:H"
	SchoolEmailRange   string = "2025-2026 Member Info!I2:I"
	PhoneNumberRange   string = "2025-2026 Member Info!J2:J"
	ShirtSizesRange    string = "2025-2026 Member Info!K2:K"
	PaidDuesRange      string = "2025-2026 Member Info!L2:L"

	LeaderRoleId  string
	OfficerRoleId string

	Officers []string

	DB *sqlx.DB

	HoursUpdateTimeout float64
	HoursLastUpdated   time.Time

	Context        context.Context
	GoogleServices *genericutils.GoogleServices
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

	LeaderRoleId = os.Getenv("LEADER_ROLE_ID")
	OfficerRoleId = os.Getenv("OFFICER_ROLE_ID")

	err := prepDatabase()
	if err != nil {
		log.Fatalf("Failed to prepare database: %v", err)
	}

	HoursTTLContender, err := strconv.Atoi(os.Getenv("HOURS_TTL"))
	if err != nil {
		log.Fatalf("Failed to convert HOURS_TTL to int: %v", err)
	}
	HoursUpdateTimeout = float64(HoursTTLContender)
	HoursLastUpdated = time.Now()

	Context = context.Background()
	GoogleServicesContender, err := genericutils.GetGoogleServices(Context, GoogleAuthKeyPath)
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
