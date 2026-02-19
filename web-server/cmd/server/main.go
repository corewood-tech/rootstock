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
	campaignflows "rootstock/web-server/flows/campaign"
	orgflows "rootstock/web-server/flows/org"
	"rootstock/web-server/global/events"
	"rootstock/web-server/global/observability"
	connecthandlers "rootstock/web-server/handlers/connect"
	campaignops "rootstock/web-server/ops/campaign"
	orgops "rootstock/web-server/ops/org"
	readingops "rootstock/web-server/ops/reading"
	"rootstock/web-server/proto/rootstock/v1/rootstockv1connect"
	"rootstock/web-server/repo/authorization"
	campaignrepo "rootstock/web-server/repo/campaign"
	eventsrepo "rootstock/web-server/repo/events"
	identityrepo "rootstock/web-server/repo/identity"
	o11yrepo "rootstock/web-server/repo/observability"
	readingrepo "rootstock/web-server/repo/reading"
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

	// Identity (JWT verification via Zitadel JWKS)
	// Zitadel resolves instances by Host header — internal requests must
	// override Host to match the external domain configured in Zitadel.
	jwksURL := fmt.Sprintf("http://%s:%d/oauth/v2/keys", cfg.Identity.Zitadel.Host, cfg.Identity.Zitadel.Port)
	jwtVerifier, err := server.NewJWTVerifier(ctx, jwksURL, cfg.Identity.Zitadel.ExternalDomain, cfg.Identity.Zitadel.Issuer)
	if err != nil {
		return fmt.Errorf("create jwt verifier: %w", err)
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
	interceptors := connect.WithInterceptors(otelInterceptor, server.AuthorizationInterceptor(jwtVerifier, authzRepo), server.BinaryOnlyInterceptor())

	// Business repos
	cRepo := campaignrepo.NewRepository(pool)
	defer cRepo.Shutdown()
	rRepo := readingrepo.NewRepository(pool)
	defer rRepo.Shutdown()

	// Identity repo (Zitadel)
	iRepo, err := identityrepo.NewRepository(ctx, cfg.Identity.Zitadel)
	if err != nil {
		return fmt.Errorf("create identity repo: %w", err)
	}
	defer iRepo.Shutdown()

	// Ops
	cOps := campaignops.NewOps(cRepo)
	rOps := readingops.NewOps(rRepo)
	oOps := orgops.NewOps(iRepo)

	// Campaign flows
	createCampaignFlow := campaignflows.NewCreateCampaignFlow(cOps)
	publishCampaignFlow := campaignflows.NewPublishCampaignFlow(cOps)
	browseCampaignsFlow := campaignflows.NewBrowseCampaignsFlow(cOps)
	campaignDashboardFlow := campaignflows.NewDashboardFlow(rOps)

	// Org flows
	createOrgFlow := orgflows.NewCreateOrgFlow(oOps)
	nestOrgFlow := orgflows.NewNestOrgFlow(oOps)
	defineRoleFlow := orgflows.NewDefineRoleFlow(oOps)
	assignRoleFlow := orgflows.NewAssignRoleFlow(oOps)
	inviteUserFlow := orgflows.NewInviteUserFlow(oOps)

	// Handlers
	healthHandler := connecthandlers.NewHealthServiceHandler()
	healthPath, healthH := rootstockv1connect.NewHealthServiceHandler(healthHandler, interceptors)

	campaignHandler := connecthandlers.NewCampaignServiceHandler(
		createCampaignFlow, publishCampaignFlow, browseCampaignsFlow, campaignDashboardFlow,
	)
	campaignPath, campaignH := rootstockv1connect.NewCampaignServiceHandler(campaignHandler, interceptors)

	orgHandler := connecthandlers.NewOrgServiceHandler(
		createOrgFlow, nestOrgFlow, defineRoleFlow, assignRoleFlow, inviteUserFlow,
	)
	orgPath, orgH := rootstockv1connect.NewOrgServiceHandler(orgHandler, interceptors)

	mux := http.NewServeMux()
	mux.Handle(healthPath, healthH)
	mux.Handle(campaignPath, campaignH)
	mux.Handle(orgPath, orgH)

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
