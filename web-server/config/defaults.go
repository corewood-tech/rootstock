package config

func defaults() Config {
	return Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Host:     "app-postgres",
				Port:     5432,
				User:     "rootstock",
				Password: "rootstock",
				DBName:   "rootstock",
				SSLMode:  "disable",
			},
		},
		Identity: IdentityConfig{
			Zitadel: ZitadelConfig{
				Host:   "zitadel",
				Port:   8080,
				Issuer:         "http://localhost:9999",
			ExternalDomain: "localhost",
			},
		},
		Authorization: AuthorizationConfig{
			OPA: OPAConfig{
				PolicyPath: "/policies",
			},
		},
		Observability: ObservabilityConfig{
			TraceExporter: "stdout",
			ServiceName:   "rootstock",
			Endpoint:      "localhost:4317",
			EnableTraces:  true,
			EnableMetrics: false,
			EnableLogs:    true,
		},
		Events: EventsConfig{
			AppName: "rootstock",
		},
		Cert: CertConfig{
			CACertPath:       "/certs/ca.crt",
			CAKeyPath:        "/certs/ca.key",
			CertLifetimeDays: 90,
		},
	}
}
