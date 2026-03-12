package database

import (
	"fmt"

	"BE_new/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsSource = "file://migrations"

func Migrate(direction string, cfg config.DatabaseConfig) error {
	m, err := migrate.New(migrationsSource, cfg.DSN)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil || databaseErr != nil {
			// Ignore close errors because the caller is already handling migration status.
		}
	}()

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		return fmt.Errorf("unsupported migration direction %q", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
