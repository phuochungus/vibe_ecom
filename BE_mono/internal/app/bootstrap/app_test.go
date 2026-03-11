package bootstrap

import (
	"testing"

	"golf-store/be-mono/internal/platform/config"
)

func TestBootstrapRequiresPostgresDSN(t *testing.T) {
	_, err := New(config.Config{ServiceName: "be-mono-test", HTTPPort: "8080", PostgresDSN: ""})
	if err == nil {
		t.Fatalf("expected bootstrap to fail when POSTGRES_DSN is empty")
	}
}
