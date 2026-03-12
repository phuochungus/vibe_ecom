package main

import (
	"log"

	"BE_new/internal/config"
	"BE_new/internal/database"
	apihttp "BE_new/internal/http"

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

	router := apihttp.NewRouter(db)
	if err := router.Run(cfg.HTTP.Address()); err != nil {
		log.Fatal(err)
	}
}
