package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"rootstock/web-server/config"
	"rootstock/web-server/global/events"
	"rootstock/web-server/global/observability"
	eventsrepo "rootstock/web-server/repo/events"
	identityrepo "rootstock/web-server/repo/identity"
	o11yrepo "rootstock/web-server/repo/observability"
	sqlconnect "rootstock/web-server/repo/sql/connect"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
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

	// Run database migrations
	if err := sqlmigrate.Run(cfg.Database.Postgres); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	logger.Info(ctx, "database migrations applied", nil)

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

	// Identity repo (Zitadel)
	iRepo, err := identityrepo.NewRepository(ctx, cfg.Identity.Zitadel)
	if err != nil {
		return fmt.Errorf("create identity repo: %w", err)
	}
	defer iRepo.Shutdown()

	// RPC server (Connect RPC, JWT auth, human traffic)
	rpcHandler, rpcCleanup, err := server.NewRPCServer(ctx, cfg, pool, iRepo)
	if err != nil {
		return fmt.Errorf("create rpc server: %w", err)
	}
	defer rpcCleanup()

	// IoT server (device HTTP, mTLS)
	iotHandler, tlsCfg, iotCleanup, err := server.NewIoTServer(cfg, pool)
	if err != nil {
		return fmt.Errorf("create iot server: %w", err)
	}
	defer iotCleanup()

	errChan := make(chan error, 2)

	// Start RPC listener (port 8080)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	rpcLis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", rpcAddr, err)
	}
	rpcServer := &http.Server{Handler: rpcHandler}
	go func() {
		logger.Info(ctx, "rpc server listening", map[string]interface{}{"addr": rpcAddr})
		if err := rpcServer.Serve(rpcLis); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("rpc serve: %w", err)
		}
	}()

	// Start IoT listener (port 8443, TLS)
	iotAddr := fmt.Sprintf("%s:8443", cfg.Server.Host)
	iotLis, err := tls.Listen("tcp", iotAddr, tlsCfg)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", iotAddr, err)
	}
	iotServer := &http.Server{Handler: iotHandler}
	go func() {
		logger.Info(ctx, "iot server listening", map[string]interface{}{"addr": iotAddr})
		if err := iotServer.Serve(iotLis); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("iot serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "shutting down servers...", nil)
		rpcServer.Close()
		iotServer.Close()
		return nil
	case err := <-errChan:
		return err
	}
}
