package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"

	"golf-store/be-mono/internal/platform/config"
	"golf-store/be-mono/internal/platform/db"
	entities "golf-store/be-mono/internal/platform/entities"
)

func main() {
	cfg := config.Load("be-mono-seed", "8080")

	database, err := db.OpenPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("[seed-demo] open database: %v", err)
	}

	if err := db.InitSchema(database); err != nil {
		log.Fatalf("[seed-demo] init schema: %v", err)
	}

	if err := db.SeedDemoData(database, cfg.PublicBaseURL); err != nil {
		log.Fatalf("[seed-demo] seed data: %v", err)
	}

	var count int64
	if err := database.Model(&entities.Product{}).Where("sku LIKE ?", "TM-%").Count(&count).Error; err != nil {
		log.Fatalf("[seed-demo] count products: %v", err)
	}

	log.Printf("[seed-demo] seeded %d TaylorMade products", count)
}
