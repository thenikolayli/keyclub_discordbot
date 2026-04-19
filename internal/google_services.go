package internal

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// struct representing an object containing all Google Services
// to be passed between functions that interact with the google APIs
type GoogleServicesType struct {
	Docs     *docs.Service
	Sheets   *sheets.Service
	Calendar *calendar.Service
}

// returns a GoogleServices struct, containing the services used to interact with the google apis
func LoadGoogleServices(ctx context.Context, keyFilePath string) (GoogleServicesType, error) {
	clientOption, err := getClientOption(keyFilePath)
	if err != nil {
		return GoogleServicesType{}, fmt.Errorf("Failed to get client option: %w", err)
	}

	docsService, err := docs.NewService(ctx, clientOption)
	if err != nil {
		return GoogleServicesType{}, fmt.Errorf("google.GetGoogleServices: %w", err)
	}
	sheetsService, err := sheets.NewService(ctx, clientOption)
	if err != nil {
		return GoogleServicesType{}, fmt.Errorf("google.GetGoogleServices: %w", err)
	}
	calendarService, err := calendar.NewService(ctx, clientOption)
	if err != nil {
		return GoogleServicesType{}, fmt.Errorf("google.GetGoogleServices: %w", err)
	}

	return GoogleServicesType{
		Docs:     docsService,
		Sheets:   sheetsService,
		Calendar: calendarService,
	}, nil
}

// uses the google_auth_key.json file to create client options
// this is used to get google services later
func getClientOption(keyFilepath string) (option.ClientOption, error) {
	if _, err := os.Stat(keyFilepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Service account key file not found: %w", err)
	}

	return option.WithAuthCredentialsFile(option.ServiceAccount, keyFilepath), nil
}
