package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"rootstock/web-server/config"
	"rootstock/web-server/global/events"
	"rootstock/web-server/global/observability"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	mqttops "rootstock/web-server/ops/mqtt"
	certrepo "rootstock/web-server/repo/cert"
	devicerepo "rootstock/web-server/repo/device"
	eventsrepo "rootstock/web-server/repo/events"
	identityrepo "rootstock/web-server/repo/identity"
	mqttrepo "rootstock/web-server/repo/mqtt"
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

	// Device repo + ops
	dRepo := devicerepo.NewRepository(pool)
	defer dRepo.Shutdown()
	dOps := deviceops.NewOps(dRepo)

	// Cert repo + ops (in-process CA, shared with RPC server for /enroll)
	crtRepo, err := certrepo.NewRepository(cfg.Cert)
	if err != nil {
		return fmt.Errorf("create cert repo: %w", err)
	}
	defer crtRepo.Shutdown()
	crtOps := certops.NewOps(crtRepo)

	// MQTT server (embedded Mochi broker, mTLS on port 8883)
	mqttServer, mqttCleanup, err := server.NewMQTTServer(cfg)
	if err != nil {
		return fmt.Errorf("create mqtt server: %w", err)
	}
	defer mqttCleanup()

	// MQTT repo + ops (wraps broker's inline client)
	mRepo := mqttrepo.NewRepository(mqttServer)
	defer mRepo.Shutdown()
	mOps := mqttops.NewOps(mRepo)

	// RPC server (Connect RPC + /enroll + /ca, returns MQTTFlows for subscription wiring)
	rpcHandler, mqttFlows, rpcCleanup, err := server.NewRPCServer(ctx, cfg, pool, iRepo, dOps, crtOps, mOps)
	if err != nil {
		return fmt.Errorf("create rpc server: %w", err)
	}
	defer rpcCleanup()

	// MQTT subscriptions (telemetry + renewal callbacks wired to flows)
	if err := server.SetupMQTTSubscriptions(ctx, mqttServer, mqttFlows); err != nil {
		return fmt.Errorf("setup mqtt subscriptions: %w", err)
	}

	errChan := make(chan error, 2)

	// Start RPC listener
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

	// Start MQTT broker
	go func() {
		logger.Info(ctx, "mqtt server starting", map[string]interface{}{"port": cfg.MQTT.Port})
		if err := mqttServer.Serve(); err != nil {
			errChan <- fmt.Errorf("mqtt serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "shutting down servers...", nil)
		rpcServer.Close()
		mqttServer.Close()
		return nil
	case err := <-errChan:
		return err
	}
}
