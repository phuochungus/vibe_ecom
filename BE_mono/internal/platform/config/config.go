package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName                  string
	HTTPPort                     string
	Env                          string
	PostgresDSN                  string
	SupabaseStorageEnabled       bool
	SupabaseURL                  string
	SupabaseServiceRoleKey       string
	SupabaseStorageBucket        string
	SupabaseStoragePublicBaseURL string
	MinIOEnabled                 bool
	MinIOEndpoint                string
	MinIOPublicBaseURL           string
	MinIOAccessKey               string
	MinIOSecretKey               string
	MinIOBucket                  string
	MinIOUseSSL                  bool
	PublicBaseURL                string
	CORSAllowedOrigins           []string
	JWTIssuer                    string
	JWTSecret                    string
	JWTAccessTTLMinutes          int
	JWTRefreshTTLMinutes         int
}

func Load(defaultServiceName string, defaultPort string) Config {
	return Config{
		ServiceName:                  getEnv("SERVICE_NAME", defaultServiceName),
		HTTPPort:                     getEnvAny([]string{"HTTP_PORT", "PORT"}, defaultPort),
		Env:                          getEnv("APP_ENV", "local"),
		PostgresDSN:                  getEnv("POSTGRES_DSN", ""),
		SupabaseStorageEnabled:       getEnvBool("SUPABASE_STORAGE_ENABLED", false),
		SupabaseURL:                  normalizePublicBaseURL(getEnv("SUPABASE_URL", "")),
		SupabaseServiceRoleKey:       getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		SupabaseStorageBucket:        getEnv("SUPABASE_STORAGE_BUCKET", "golf-store"),
		SupabaseStoragePublicBaseURL: normalizePublicBaseURL(getEnv("SUPABASE_STORAGE_PUBLIC_BASE_URL", "")),
		MinIOEnabled:                 getEnvBool("MINIO_ENABLED", false),
		MinIOEndpoint:                getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		MinIOPublicBaseURL:           normalizePublicBaseURL(getEnv("MINIO_PUBLIC_BASE_URL", "http://127.0.0.1:9000")),
		MinIOAccessKey:               getEnv("MINIO_ACCESS_KEY", ""),
		MinIOSecretKey:               getEnv("MINIO_SECRET_KEY", ""),
		MinIOBucket:                  getEnv("MINIO_BUCKET", "golf-store"),
		MinIOUseSSL:                  getEnvBool("MINIO_USE_SSL", false),
		PublicBaseURL:                normalizePublicBaseURL(getEnv("PUBLIC_BASE_URL", "http://127.0.0.1:"+defaultPort)),
		CORSAllowedOrigins:           getEnvCSV("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3001", "http://127.0.0.1:3001", "http://localhost:4173", "http://127.0.0.1:4173"}),
		JWTIssuer:                    getEnv("JWT_ISSUER", defaultServiceName),
		JWTSecret:                    getEnv("JWT_SECRET", "dev_jwt_secret_change_me"),
		JWTAccessTTLMinutes:          getEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
		JWTRefreshTTLMinutes:         getEnvInt("JWT_REFRESH_TTL_MINUTES", 10080),
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

func getEnvBool(key string, fallback bool) bool {
	v := strings.ToLower(strings.TrimSpace(getEnv(key, "")))
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
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
