package main

import (
	"log"

	"BE_new/internal/config"
	"BE_new/internal/database"
	"BE_new/internal/routes"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Open(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	router := routes.NewRouter(db)
	if err := router.Run(cfg.HTTP.Address()); err != nil {
		log.Fatal(err)
	}
}
