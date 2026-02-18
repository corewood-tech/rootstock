package observability

import (
	"context"

	o11yrepo "rootstock/web-server/repo/observability"
)

type getTracerReq struct {
	name string
	resp chan o11yrepo.Tracer
}

type getMeterReq struct {
	name string
	resp chan o11yrepo.Meter
}

type getLoggerReq struct {
	name string
	resp chan o11yrepo.Logger
}

type shutdownReq struct {
	ctx  context.Context
	resp chan error
}

var (
	tracerCh   = make(chan getTracerReq)
	meterCh    = make(chan getMeterReq)
	loggerCh   = make(chan getLoggerReq)
	shutdownCh = make(chan shutdownReq)
)

// Initialize starts the manager goroutine that delegates to the provided repository.
func Initialize(repo o11yrepo.Repository) {
	go manage(repo)
}

// manage owns the repo reference. All access goes through channels.
func manage(repo o11yrepo.Repository) {
	for {
		select {
		case req := <-tracerCh:
			req.resp <- repo.GetTracer(req.name)

		case req := <-meterCh:
			req.resp <- repo.GetMeter(req.name)

		case req := <-loggerCh:
			req.resp <- repo.GetLogger(req.name)

		case req := <-shutdownCh:
			req.resp <- repo.Shutdown(req.ctx)
			return
		}
	}
}

// Shutdown gracefully shuts down all providers.
func Shutdown(ctx context.Context) error {
	resp := make(chan error, 1)
	shutdownCh <- shutdownReq{ctx: ctx, resp: resp}
	return <-resp
}

// GetTracer returns a named tracer from the global provider.
func GetTracer(name string) o11yrepo.Tracer {
	resp := make(chan o11yrepo.Tracer, 1)
	tracerCh <- getTracerReq{name: name, resp: resp}
	return <-resp
}

// GetMeter returns a named meter from the global provider.
func GetMeter(name string) o11yrepo.Meter {
	resp := make(chan o11yrepo.Meter, 1)
	meterCh <- getMeterReq{name: name, resp: resp}
	return <-resp
}

// GetLogger returns a Logger backed by the configured provider.
func GetLogger(name string) o11yrepo.Logger {
	resp := make(chan o11yrepo.Logger, 1)
	loggerCh <- getLoggerReq{name: name, resp: resp}
	return <-resp
}
