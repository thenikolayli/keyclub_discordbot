package config

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"log"
	"os"
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

	NamesRange         string = "2025-2026 Members!A2:A"
	NicknamesRange     string = "2025-2026 Members!B2:B"
	AllHoursRange      string = "2025-2026 Members!C2:C"
	TermHoursRange     string = "2025-2026 Members!D2:D"
	GradYearRange      string = "2025-2026 Members!E2:E"
	ClassRange         string = "2025-2026 Members!F2:F"
	StrikesRange       string = "2025-2026 Members!G2:G"
	PersonalEmailRange string = "2025-2026 Members!H2:H"
	SchoolEmailRange   string = "2025-2026 Members!I2:I"
	PhoneNumberRange   string = "2025-2026 Members!J2:J"
	ShirtSizesRange    string = "2025-2026 Members!K2:K"
	PaidDuesRange      string = "2025-2026 Members!L2:L"

	LeaderRoleId  string
	OfficerRoleId string

	Officers []string

	DB *sqlx.DB

	HoursUpdateTimeout float64 = 3600
	HoursLastUpdated   time.Time

	DefaultRankTopN int = 5

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
