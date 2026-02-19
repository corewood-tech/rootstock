package migrate

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
)

var testCfg = config.PostgresConfig{
	Host:     "app-postgres",
	Port:     5432,
	User:     "rootstock",
	Password: "rootstock",
	DBName:   "rootstock",
	SSLMode:  "disable",
}

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		testCfg.User, testCfg.Password, testCfg.Host, testCfg.Port, testCfg.DBName, testCfg.SSLMode,
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestRunAppliesMigrations(t *testing.T) {
	if err := Run(testCfg); err != nil {
		t.Fatalf("Run() first call: %v", err)
	}

	// Idempotent
	if err := Run(testCfg); err != nil {
		t.Fatalf("Run() second call (idempotent): %v", err)
	}
}

func TestMigrationsCreateExpectedTables(t *testing.T) {
	if err := Run(testCfg); err != nil {
		t.Fatalf("Run(): %v", err)
	}

	pool := testPool(t)
	ctx := context.Background()

	expected := []string{
		"campaigns",
		"campaign_parameters",
		"campaign_regions",
		"campaign_eligibility",
		"devices",
		"enrollment_codes",
		"device_campaigns",
		"readings",
		"scores",
		"badges",
		"sweepstakes_entries",
	}

	for _, table := range expected {
		var exists bool
		err := pool.QueryRow(ctx,
			"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)",
			table,
		).Scan(&exists)
		if err != nil {
			t.Errorf("query for table %s: %v", table, err)
		} else if !exists {
			t.Errorf("table %s does not exist", table)
		}
	}
}
