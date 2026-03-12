package main

import (
	"fmt"
	"log"
	"os"

	"BE_new/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	direction, err := migrationDirection(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	dbConfig, err := config.LoadDatabase()
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New("file://db/migrations", dbConfig.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("close migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("close migration database: %v", databaseErr)
		}
	}()

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("unsupported migration direction %q", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func migrationDirection(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("usage: go run ./cmd/migrate [up|down]")
	}

	switch args[0] {
	case "up", "down":
		return args[0], nil
	default:
		return "", fmt.Errorf("usage: go run ./cmd/migrate [up|down]")
	}
}
