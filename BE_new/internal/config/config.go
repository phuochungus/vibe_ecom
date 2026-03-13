package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
}

type HTTPConfig struct {
	Port string
}

func (c HTTPConfig) Address() string {
	return ":" + c.Port
}

type DatabaseConfig struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSL_MODE string
	DB_TIMEZONE string
	DSN         string
}

func Load() (Config, error) {
	dbConfig, err := LoadDatabase()
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTP: HTTPConfig{
			Port: os.Getenv("APP_PORT"),
		},
		Database: dbConfig,
	}, nil
}

func LoadDatabase() (DatabaseConfig, error) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_SSL_MODE := os.Getenv("DB_SSL_MODE")
	DB_TIMEZONE := os.Getenv("DB_TIMEZONE")
	if DB_HOST == "" || DB_PORT == "" || DB_USER == "" || DB_PASSWORD == "" || DB_NAME == "" {
		return DatabaseConfig{}, fmt.Errorf("database configuration is incomplete")
	}
	DSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT,
		DB_SSL_MODE,
		DB_TIMEZONE,
	)

	return DatabaseConfig{
		DB_HOST:     DB_HOST,
		DB_PORT:     DB_PORT,
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_NAME:     DB_NAME,
		DB_SSL_MODE: DB_SSL_MODE,
		DB_TIMEZONE: DB_TIMEZONE,
		DSN:         DSN,
	}, nil
}

func (c DatabaseConfig) MigrationURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		url.QueryEscape(c.DB_USER),
		url.QueryEscape(c.DB_PASSWORD),
		c.DB_HOST,
		c.DB_PORT,
		c.DB_NAME,
		c.DB_SSL_MODE,
		url.QueryEscape(c.DB_TIMEZONE),
	)
}
