package config

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	SpreadsheetID     string
	GoogleAuthKeyPath string

	SheetName      string
	NamesRange     string
	NicknamesRange string
	TermHoursRange string
	AllHoursRange  string
	GradYearRange  string
	StrikesRange   string

	DB *sqlx.DB

	HoursTTL         float64
	HoursLastUpdated time.Time

	Context            context.Context
	GoogleServices *genericutils.GoogleServices
)

const memberSchema = `
CREATE TABLE IF NOT EXISTS members (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	nickname TEXT,
	term_hours FLOAT NOT NULL,
	all_hours FLOAT NOT NULL,
	shirt_size TEXT,
	paid_dues BOOLEAN,
	grad_year INTEGER,
	strikes INTEGER
)
`

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	DBContender, err := sqlx.Connect("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	DB = DBContender
	DB.MustExec(memberSchema) // creates the members table

	SpreadsheetID = os.Getenv("HOURS_SPREADSHEET_ID")
	GoogleAuthKeyPath = os.Getenv("GOOGLE_AUTH_KEY_PATH")

	SheetName = os.Getenv("SHEET_NAME")
	NamesRange = fmt.Sprintf("%v!%v", SheetName, os.Getenv("NAMES_RANGE"))
	NicknamesRange = fmt.Sprintf("%v!%v", SheetName, os.Getenv("NICKNAMES_RANGE"))
	TermHoursRange = fmt.Sprintf("%v!%v", SheetName, os.Getenv("TERM_HOURS_RANGE"))
	AllHoursRange = fmt.Sprintf("%v!%v", SheetName, os.Getenv("ALL_HOURS_RANGE"))
	GradYearRange = fmt.Sprintf("%v!%v", SheetName, os.Getenv("GRAD_YEAR_RANGE"))
	StrikesRange = fmt.Sprintf("%v!%v", SheetName, os.Getenv("STRIKES_RANGE"))

	HoursTTLContender, err := strconv.Atoi(os.Getenv("HOURS_TTL"))
	if err != nil {
		log.Fatalf("Failed to convert HOURS_TTL to int: %v", err)
	}
	HoursTTL = float64(HoursTTLContender)

	Context = context.Background()
	GoogleServicesContender, err := genericutils.GetGoogleServices(Context, GoogleAuthKeyPath)
	if err != nil {
		log.Fatalf("Issue getting Google services: %v", err)
	}
	GoogleServices = GoogleServicesContender
}
