package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServiceName          string
	HTTPPort             string
	Env                  string
	MySQLDSN             string
	JWTIssuer            string
	JWTSecret            string
	JWTAccessTTLMinutes  int
	JWTRefreshTTLMinutes int
}

func Load(defaultServiceName string, defaultPort string) Config {
	return Config{
		ServiceName:          getEnv("SERVICE_NAME", defaultServiceName),
		HTTPPort:             getEnv("HTTP_PORT", defaultPort),
		Env:                  getEnv("APP_ENV", "local"),
		MySQLDSN:             getEnv("MYSQL_DSN", ""),
		JWTIssuer:            getEnv("JWT_ISSUER", defaultServiceName),
		JWTSecret:            getEnv("JWT_SECRET", "dev_jwt_secret_change_me"),
		JWTAccessTTLMinutes:  getEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
		JWTRefreshTTLMinutes: getEnvInt("JWT_REFRESH_TTL_MINUTES", 10080),
	}
}

func getEnv(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := getEnv(key, "")
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
