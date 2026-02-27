package bootstrap

import (
	"testing"

	"golf-store/be-mono/internal/platform/config"
)

func TestBootstrapRequiresMySQLDSN(t *testing.T) {
	_, err := New(config.Config{ServiceName: "be-mono-test", HTTPPort: "8080", MySQLDSN: ""})
	if err == nil {
		t.Fatalf("expected bootstrap to fail when MYSQL_DSN is empty")
	}
}
