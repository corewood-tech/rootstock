package migrate

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"rootstock/web-server/config"
)

//go:embed migrations/*.sql
var fs embed.FS

// Run applies all pending migrations. It embeds the SQL files from the
// migrations directory so there are no runtime filesystem dependencies.
func Run(cfg config.PostgresConfig) error {
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	dsn := fmt.Sprintf(
		"pgx5://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
