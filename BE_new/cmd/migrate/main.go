package main

import (
	"fmt"
	"log"
	"os"

	"BE_new/internal/config"
	"BE_new/internal/database"

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

	if err := database.Migrate(direction, dbConfig); err != nil {
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
