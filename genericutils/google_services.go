package genericutils

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// returns a GoogleServices struct, containing the services used to interact with the google apis
func GetGoogleServices(ctx context.Context, keyFilePath string) (*GoogleServices, error) {
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

	return &GoogleServices{
		Docs:   docsService,
		Sheets: sheetsService,
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
