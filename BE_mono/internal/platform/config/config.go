package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName          string
	HTTPPort             string
	Env                  string
	PostgresDSN          string
	PublicBaseURL        string
	CORSAllowedOrigins   []string
	JWTIssuer            string
	JWTSecret            string
	JWTAccessTTLMinutes  int
	JWTRefreshTTLMinutes int
}

func Load(defaultServiceName string, defaultPort string) Config {
	return Config{
		ServiceName:          getEnv("SERVICE_NAME", defaultServiceName),
		HTTPPort:             getEnvAny([]string{"HTTP_PORT", "PORT"}, defaultPort),
		Env:                  getEnv("APP_ENV", "local"),
		PostgresDSN:          getEnv("POSTGRES_DSN", ""),
		PublicBaseURL:        normalizePublicBaseURL(getEnv("PUBLIC_BASE_URL", "http://127.0.0.1:"+defaultPort)),
		CORSAllowedOrigins:   getEnvCSV("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3001", "http://127.0.0.1:3001", "http://localhost:4173", "http://127.0.0.1:4173"}),
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

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		if value := getEnv(key, ""); value != "" {
			return value
		}
	}
	return fallback
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

func getEnvCSV(key string, fallback []string) []string {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback
	}

	values := make([]string, 0)
	current := ""
	for _, r := range raw {
		if r == ',' {
			if current != "" {
				values = append(values, current)
			}
			current = ""
			continue
		}
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		current += string(r)
	}
	if current != "" {
		values = append(values, current)
	}
	if len(values) == 0 {
		return fallback
	}
	return values
}

func normalizePublicBaseURL(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}
	return "https://" + trimmed
}
