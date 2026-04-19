package internal

import (
	"context"
	"fmt"
	"keyclubDiscordBot/genericutils"
	"time"

	"github.com/jmoiron/sqlx"
)

type App struct {
	Config         Config
	GoogleServices GoogleServicesType
	DB             *sqlx.DB
	MemberSync     genericutils.SyncState
	EventSync      genericutils.SyncState
	// ShutdownCtx is cancelled on process shutdown (e.g. SIGINT/SIGTERM). Handlers should derive
	// per-request timeouts with context.WithTimeout(ShutdownCtx, ...).
	ShutdownCtx context.Context
}

func NewApp(ctx context.Context) (*App, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %w", err)
	}
	googleServices, err := LoadGoogleServices(ctx, config.GoogleAuthKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load google services: %w", err)
	}
	db, err := LoadDatabase()
	if err != nil {
		return nil, fmt.Errorf("Failed to load database: %w", err)
	}

	return &App{
		Config:         config,
		GoogleServices: googleServices,
		DB:             db,
		MemberSync:     genericutils.SyncState{UpdateTimeout: config.MemberSyncTimeout, LastUpdated: time.Date(2026, time.January, 1, 01, 01, 0, 0, time.UTC)},
		EventSync:      genericutils.SyncState{UpdateTimeout: config.EventSyncTimeout, LastUpdated: time.Date(2026, time.January, 1, 01, 01, 0, 0, time.UTC)},
		ShutdownCtx:    ctx,
	}, nil
}
