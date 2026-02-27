package config

import "os"

type Config struct {
	ServiceName string
	HTTPPort    string
	Env         string
	MySQLDSN    string
}

func Load(defaultServiceName string, defaultPort string) Config {
	return Config{
		ServiceName: getEnv("SERVICE_NAME", defaultServiceName),
		HTTPPort:    getEnv("HTTP_PORT", defaultPort),
		Env:         getEnv("APP_ENV", "local"),
		MySQLDSN:    getEnv("MYSQL_DSN", ""),
	}
}

func getEnv(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	return v
}
