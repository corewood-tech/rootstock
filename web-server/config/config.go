package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type Config struct {
	Server        ServerConfig        `koanf:"server"`
	Database      DatabaseConfig      `koanf:"database"`
	Identity      IdentityConfig      `koanf:"identity"`
	Authorization AuthorizationConfig `koanf:"authorization"`
	Observability ObservabilityConfig `koanf:"observability"`
	Events        EventsConfig        `koanf:"events"`
}

type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `koanf:"postgres"`
}

type PostgresConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	DBName   string `koanf:"dbname"`
	SSLMode  string `koanf:"sslmode"`
}

type IdentityConfig struct {
	Zitadel ZitadelConfig `koanf:"zitadel"`
}

type ZitadelConfig struct {
	Host           string `koanf:"host"`
	Port           int    `koanf:"port"`
	Issuer         string `koanf:"issuer"`
	ExternalDomain string `koanf:"external_domain"`
	ServiceUserPAT string `koanf:"service_user_pat"`
}

type AuthorizationConfig struct {
	OPA OPAConfig `koanf:"opa"`
}

type OPAConfig struct {
	PolicyPath string `koanf:"policy_path"`
}

type ObservabilityConfig struct {
	TraceExporter string `koanf:"trace_exporter"`
	ServiceName   string `koanf:"service_name"`
	Endpoint      string `koanf:"endpoint"`
	EnableTraces  bool   `koanf:"enable_traces"`
	EnableMetrics bool   `koanf:"enable_metrics"`
	EnableLogs    bool   `koanf:"enable_logs"`
}

type EventsConfig struct {
	AppName string `koanf:"app_name"`
}

// Load builds the config by layering: defaults → YAML file → env vars → CLI flags.
func Load(configPath string, flags *pflag.FlagSet) (*Config, error) {
	k := koanf.New(".")

	// 1. Struct defaults
	if err := k.Load(structs.Provider(defaults(), "koanf"), nil); err != nil {
		return nil, fmt.Errorf("load defaults: %w", err)
	}

	// 2. YAML config file (optional)
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("load config file %s: %w", configPath, err)
		}
	}

	// 3. Environment variables: ROOTSTOCK_SERVER_PORT → server.port
	if err := k.Load(env.Provider("ROOTSTOCK_", ".", func(s string) string {
		return strings.Replace(
			strings.ToLower(strings.TrimPrefix(s, "ROOTSTOCK_")),
			"_", ".", -1,
		)
	}), nil); err != nil {
		return nil, fmt.Errorf("load env vars: %w", err)
	}

	// 4. CLI flags (optional)
	if flags != nil {
		if err := k.Load(posflag.Provider(flags, ".", k), nil); err != nil {
			return nil, fmt.Errorf("load flags: %w", err)
		}
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
