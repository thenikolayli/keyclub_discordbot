package config

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	SpreadsheetID     string
	GoogleAuthKeyPath string

	NamesRange     string
	NicknamesRange string
	TermHoursRange string
	AllHoursRange  string
	GradYearRange  string
	StrikesRange   string

	DB *sqlx.DB

	HoursUpdateTimeout float64
	HoursLastUpdated   time.Time

	Context        context.Context
	GoogleServices *genericutils.GoogleServices
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	err := prepDatabase()
	if err != nil {
		log.Fatalf("Failed to prepare database: %v", err)
	}

	SpreadsheetID = os.Getenv("HOURS_SPREADSHEET_ID")
	GoogleAuthKeyPath = os.Getenv("GOOGLE_AUTH_KEY_PATH")

	NamesRange = os.Getenv("NAMES_RANGE")
	NicknamesRange = os.Getenv("NICKNAMES_RANGE")
	TermHoursRange = os.Getenv("TERM_HOURS_RANGE")
	AllHoursRange = os.Getenv("ALL_HOURS_RANGE")
	GradYearRange = os.Getenv("GRAD_YEAR_RANGE")
	StrikesRange = os.Getenv("STRIKES_RANGE")

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
}

// prepares the database and runs migrations
func prepDatabase() error {
	DBContender, err := sqlx.Connect("sqlite3", "db.sqlite3")
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}
	DB = DBContender

	migrations, migrationErr := migrate.New(
		"file://migrations",
		"sqlite3://db.sqlite3",
	)
	if migrationErr != nil {
		return fmt.Errorf("Failed to initialize migrations: %w", migrationErr)
	}
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}
	return nil
}
