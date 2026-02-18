package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"

	"rootstock/web-server/config"
	"rootstock/web-server/global/events"
	"rootstock/web-server/global/observability"
	connecthandlers "rootstock/web-server/handlers/connect"
	"rootstock/web-server/proto/rootstock/v1/rootstockv1connect"
	"rootstock/web-server/repo/authorization"
	eventsrepo "rootstock/web-server/repo/events"
	o11yrepo "rootstock/web-server/repo/observability"
	sqlconnect "rootstock/web-server/repo/sql/connect"
	"rootstock/web-server/server"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load("config.yaml", nil)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Observability repo — first thing after config
	o11y, err := o11yrepo.NewOTelRepository(ctx, cfg.Observability)
	if err != nil {
		return fmt.Errorf("create observability repo: %w", err)
	}
	observability.Initialize(o11y)
	defer observability.Shutdown(ctx)

	logger := observability.GetLogger("main")

	// Database pool
	pool, err := sqlconnect.OpenPostgres(ctx, cfg.Database.Postgres)
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	defer pool.Close()

	// Events repo — uses injected pool
	evts, err := eventsrepo.NewDBOSRepository(ctx, pool, cfg.Events.AppName)
	if err != nil {
		return fmt.Errorf("create events repo: %w", err)
	}
	events.Initialize(evts)
	defer events.Shutdown()

	// Runtime metrics (goroutines, memory, GC)
	if err := runtime.Start(); err != nil {
		return fmt.Errorf("start runtime metrics: %w", err)
	}

	// Authorization (OPA)
	authzRepo := authorization.NewOPARepository()
	if err := authzRepo.Recompile(ctx); err != nil {
		return fmt.Errorf("compile authorization policy: %w", err)
	}

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return fmt.Errorf("create otel interceptor: %w", err)
	}
	interceptors := connect.WithInterceptors(otelInterceptor, server.AuthorizationInterceptor(authzRepo), server.BinaryOnlyInterceptor())

	healthHandler := connecthandlers.NewHealthServiceHandler()
	path, handler := rootstockv1connect.NewHealthServiceHandler(healthHandler, interceptors)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	httpServer := &http.Server{Handler: mux}

	errChan := make(chan error, 1)
	go func() {
		logger.Info(ctx, "server listening", map[string]interface{}{"addr": addr})
		if err := httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "shutting down server...", nil)
		return httpServer.Close()
	case err := <-errChan:
		return err
	}
}
