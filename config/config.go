package config

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

	HoursTTL         int64
	HoursLastUpdated int64 = 0

	ctx            context.Context
	GoogleServices *genericutils.GoogleServices
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

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
	HoursTTL = int64(HoursTTLContender)

	ctx = context.Background()
	GoogleServicesContender, err := genericutils.GetGoogleServices(ctx, GoogleAuthKeyPath)
	if err != nil {
		log.Fatalf("Issue getting Google services: %v", err)
	}
	GoogleServices = GoogleServicesContender
}
