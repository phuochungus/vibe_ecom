package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	DB_HOST := "127.0.0.1"
	DB_PORT := "5433"
	DB_USER := "golf"
	DB_PASSWORD := "golf"
	DB_NAME := "golf_store_mono"
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://"+DB_USER+":"+DB_PASSWORD+"@"+DB_HOST+":"+DB_PORT+"/"+DB_NAME+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
}
