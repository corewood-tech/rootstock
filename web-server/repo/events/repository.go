package events

import (
	"context"
	"fmt"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type getContextReq struct {
	resp chan WorkflowContext
}

type shutdownReq struct {
	resp chan struct{}
}

type dbosRepo struct {
	getContextCh chan getContextReq
	shutdownCh   chan shutdownReq
}

// NewDBOSRepository creates an events repository backed by DBOS.
// It initializes and launches the DBOS runtime. The caller owns the pool lifecycle.
func NewDBOSRepository(ctx context.Context, pool *pgxpool.Pool, appName string) (Repository, error) {
	cfg := dbos.Config{
		AppName:      appName,
		SystemDBPool: pool,
	}

	dbosCtx, err := dbos.NewDBOSContext(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create DBOS context: %w", err)
	}

	if err := dbos.Launch(dbosCtx); err != nil {
		return nil, fmt.Errorf("launch DBOS: %w", err)
	}

	r := &dbosRepo{
		getContextCh: make(chan getContextReq),
		shutdownCh:   make(chan shutdownReq),
	}
	go r.manage(dbosCtx)
	return r, nil
}

func (r *dbosRepo) manage(dbosCtx dbos.DBOSContext) {
	for {
		select {
		case req := <-r.getContextCh:
			req.resp <- dbosCtx

		case req := <-r.shutdownCh:
			dbos.Shutdown(dbosCtx, 30*time.Second)
			req.resp <- struct{}{}
			return
		}
	}
}

func (r *dbosRepo) GetContext() WorkflowContext {
	resp := make(chan WorkflowContext, 1)
	r.getContextCh <- getContextReq{resp: resp}
	return <-resp
}

func (r *dbosRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}
