package internal

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

// prepares the database and runs migrations
func LoadDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", "db.sqlite3")
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %w", err)
	}

	migrations, migrationErr := migrate.New(
		"file://migrations",
		"sqlite://db.sqlite3",
	)
	if migrationErr != nil {
		return nil, fmt.Errorf("Failed to initialize migrations: %w", migrationErr)
	}
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Failed to run migrations: %w", err)
	}
	return db, nil
}
