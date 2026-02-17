package events

import (
	"context"
	"fmt"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type getContextReq struct {
	resp chan dbos.DBOSContext
}

type shutdownReq struct {
	resp chan struct{}
}

var (
	getContextCh = make(chan getContextReq)
	shutdownCh   = make(chan shutdownReq)
)

// Initialize creates a DBOS context using the provided pool and launches it.
// The caller owns the pool lifecycle — Shutdown does not close it.
func Initialize(ctx context.Context, pool *pgxpool.Pool, appName string) error {
	cfg := dbos.Config{
		AppName:      appName,
		SystemDBPool: pool,
	}

	dbosCtx, err := dbos.NewDBOSContext(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create DBOS context: %w", err)
	}

	if err := dbos.Launch(dbosCtx); err != nil {
		return fmt.Errorf("launch DBOS: %w", err)
	}

	go manage(dbosCtx)
	return nil
}

// manage owns the DBOS context. All access goes through channels.
func manage(dbosCtx dbos.DBOSContext) {
	for {
		select {
		case req := <-getContextCh:
			req.resp <- dbosCtx

		case req := <-shutdownCh:
			dbos.Shutdown(dbosCtx, 30*time.Second)
			req.resp <- struct{}{}
			return
		}
	}
}

// Shutdown gracefully shuts down the DBOS runtime.
// It does not close the database pool.
func Shutdown() {
	resp := make(chan struct{}, 1)
	shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// GetContext returns the DBOS context for workflow execution.
func GetContext() dbos.DBOSContext {
	resp := make(chan dbos.DBOSContext, 1)
	getContextCh <- getContextReq{resp: resp}
	return <-resp
}
