package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("", nil)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.Database.Postgres.Host != "app-postgres" {
		t.Errorf("Database.Postgres.Host = %q, want %q", cfg.Database.Postgres.Host, "app-postgres")
	}
	if cfg.Observability.ServiceName != "rootstock" {
		t.Errorf("Observability.ServiceName = %q, want %q", cfg.Observability.ServiceName, "rootstock")
	}
}

func TestLoadFileOverridesDefaults(t *testing.T) {
	content := []byte(`
server:
  port: 9090
database:
  postgres:
    host: db-host
`)
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := Load(path, nil)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 9090)
	}
	if cfg.Database.Postgres.Host != "db-host" {
		t.Errorf("Database.Postgres.Host = %q, want %q", cfg.Database.Postgres.Host, "db-host")
	}
	// Defaults should still be present for unset fields
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want default %q", cfg.Server.Host, "0.0.0.0")
	}
}

func TestLoadEnvOverridesFile(t *testing.T) {
	content := []byte(`
server:
  port: 9090
`)
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	t.Setenv("ROOTSTOCK_SERVER_PORT", "7070")
	t.Setenv("ROOTSTOCK_DATABASE_POSTGRES_HOST", "env-db-host")

	cfg, err := Load(path, nil)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 7070 {
		t.Errorf("Server.Port = %d, want %d (env override)", cfg.Server.Port, 7070)
	}
	if cfg.Database.Postgres.Host != "env-db-host" {
		t.Errorf("Database.Postgres.Host = %q, want %q (env override)", cfg.Database.Postgres.Host, "env-db-host")
	}
}

func TestLoadFlagsOverrideEnv(t *testing.T) {
	t.Setenv("ROOTSTOCK_SERVER_PORT", "7070")

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Int("server.port", 0, "server port")
	flags.Parse([]string{"--server.port=3000"})

	cfg, err := Load("", flags)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 3000 {
		t.Errorf("Server.Port = %d, want %d (flag override)", cfg.Server.Port, 3000)
	}
}

func TestLoadFullPrecedence(t *testing.T) {
	// File sets port to 9090
	content := []byte(`
server:
  port: 9090
  host: file-host
`)
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	// Env overrides host
	t.Setenv("ROOTSTOCK_SERVER_HOST", "env-host")

	// Flag overrides port
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Int("server.port", 0, "server port")
	flags.Parse([]string{"--server.port=5555"})

	cfg, err := Load(path, flags)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Flag wins over file for port
	if cfg.Server.Port != 5555 {
		t.Errorf("Server.Port = %d, want %d (flag wins)", cfg.Server.Port, 5555)
	}
	// Env wins over file for host
	if cfg.Server.Host != "env-host" {
		t.Errorf("Server.Host = %q, want %q (env wins)", cfg.Server.Host, "env-host")
	}
	// Default still present for untouched fields
	if cfg.Database.Postgres.SSLMode != "disable" {
		t.Errorf("Database.Postgres.SSLMode = %q, want default %q", cfg.Database.Postgres.SSLMode, "disable")
	}
}
