package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	mu          sync.RWMutex
	dbosCtx     dbos.DBOSContext
	initialized bool
)

// Initialize creates a DBOS context using the provided pool and launches it.
// The caller owns the pool lifecycle — Shutdown does not close it.
func Initialize(ctx context.Context, pool *pgxpool.Pool, appName string) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	cfg := dbos.Config{
		AppName:      appName,
		SystemDBPool: pool,
	}

	var err error
	dbosCtx, err = dbos.NewDBOSContext(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create DBOS context: %w", err)
	}

	if err := dbos.Launch(dbosCtx); err != nil {
		return fmt.Errorf("launch DBOS: %w", err)
	}

	initialized = true
	return nil
}

// Shutdown gracefully shuts down the DBOS runtime.
// It does not close the database pool.
func Shutdown() {
	mu.Lock()
	defer mu.Unlock()

	if !initialized {
		return
	}

	dbos.Shutdown(dbosCtx, 30*time.Second)
	initialized = false
}

// GetContext returns the DBOS context for workflow execution.
func GetContext() dbos.DBOSContext {
	mu.RLock()
	defer mu.RUnlock()

	return dbosCtx
}
