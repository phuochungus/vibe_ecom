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
	DSN string
}

func Load() (Config, error) {
	dbConfig, err := LoadDatabase()
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTP: HTTPConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: dbConfig,
	}, nil
}

func LoadDatabase() (DatabaseConfig, error) {
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		return DatabaseConfig{DSN: dsn}, nil
	}

	user := getEnv("DB_USER", "golf")
	password := getEnv("DB_PASSWORD", "golf")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5433")
	name := getEnv("DB_NAME", "golf_store_mono")
	sslMode := getEnv("DB_SSLMODE", "disable")

	if name == "" {
		return DatabaseConfig{}, fmt.Errorf("DB_NAME environment variable is required")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		name,
		sslMode,
	)

	return DatabaseConfig{DSN: dsn}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
