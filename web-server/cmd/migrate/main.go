package main

import (
	"fmt"
	"os"

	"rootstock/web-server/config"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load("config.yaml", nil)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := sqlmigrate.Run(cfg.Database.Postgres); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	fmt.Println("migrations applied successfully")
	return nil
}
