package server

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	campaignflows "rootstock/web-server/flows/campaign"
	deviceflows "rootstock/web-server/flows/device"
	orgflows "rootstock/web-server/flows/org"
	scoreflows "rootstock/web-server/flows/score"
	connecthandlers "rootstock/web-server/handlers/connect"
	campaignops "rootstock/web-server/ops/campaign"
	deviceops "rootstock/web-server/ops/device"
	orgops "rootstock/web-server/ops/org"
	readingops "rootstock/web-server/ops/reading"
	scoreops "rootstock/web-server/ops/score"
	"rootstock/web-server/proto/rootstock/v1/rootstockv1connect"
	"rootstock/web-server/repo/authorization"
	campaignrepo "rootstock/web-server/repo/campaign"
	devicerepo "rootstock/web-server/repo/device"
	identityrepo "rootstock/web-server/repo/identity"
	readingrepo "rootstock/web-server/repo/reading"
	scorerepo "rootstock/web-server/repo/score"
)

// NewRPCServer wires repos → ops → flows → Connect RPC handlers and returns
// an http.Handler + a shutdown function.
func NewRPCServer(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, iRepo identityrepo.Repository) (http.Handler, func(), error) {
	// JWT verification
	jwksURL := fmt.Sprintf("http://%s:%d/oauth/v2/keys", cfg.Identity.Zitadel.Host, cfg.Identity.Zitadel.Port)
	jwtVerifier, err := NewJWTVerifier(ctx, jwksURL, cfg.Identity.Zitadel.ExternalDomain, cfg.Identity.Zitadel.Issuer)
	if err != nil {
		return nil, nil, fmt.Errorf("create jwt verifier: %w", err)
	}

	// Authorization (OPA)
	authzRepo := authorization.NewOPARepository()
	if err := authzRepo.Recompile(ctx); err != nil {
		return nil, nil, fmt.Errorf("compile authorization policy: %w", err)
	}

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, nil, fmt.Errorf("create otel interceptor: %w", err)
	}
	interceptors := connect.WithInterceptors(otelInterceptor, AuthorizationInterceptor(jwtVerifier, authzRepo), BinaryOnlyInterceptor())

	// Business repos
	cRepo := campaignrepo.NewRepository(pool)
	dRepo := devicerepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)
	sRepo := scorerepo.NewRepository(pool)

	// Ops
	cOps := campaignops.NewOps(cRepo)
	dOps := deviceops.NewOps(dRepo)
	rOps := readingops.NewOps(rRepo)
	oOps := orgops.NewOps(iRepo)
	sOps := scoreops.NewOps(sRepo)

	// Campaign flows
	createCampaignFlow := campaignflows.NewCreateCampaignFlow(cOps)
	publishCampaignFlow := campaignflows.NewPublishCampaignFlow(cOps)
	browseCampaignsFlow := campaignflows.NewBrowseCampaignsFlow(cOps)
	campaignDashboardFlow := campaignflows.NewDashboardFlow(rOps)

	// Device flows
	getDeviceFlow := deviceflows.NewGetDeviceFlow(dOps)
	revokeDeviceFlow := deviceflows.NewRevokeDeviceFlow(dOps)
	reinstateDeviceFlow := deviceflows.NewReinstateDeviceFlow(dOps)

	// Score flows
	getContributionFlow := scoreflows.NewGetContributionFlow(sOps)

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

	scoreHandler := connecthandlers.NewScoreServiceHandler(getContributionFlow)
	scorePath, scoreH := rootstockv1connect.NewScoreServiceHandler(scoreHandler, interceptors)

	deviceHandler := connecthandlers.NewDeviceServiceHandler(getDeviceFlow, revokeDeviceFlow, reinstateDeviceFlow)
	devicePath, deviceH := rootstockv1connect.NewDeviceServiceHandler(deviceHandler, interceptors)

	mux := http.NewServeMux()
	mux.Handle(healthPath, healthH)
	mux.Handle(campaignPath, campaignH)
	mux.Handle(orgPath, orgH)
	mux.Handle(scorePath, scoreH)
	mux.Handle(devicePath, deviceH)

	shutdown := func() {
		cRepo.Shutdown()
		dRepo.Shutdown()
		rRepo.Shutdown()
		sRepo.Shutdown()
	}

	return mux, shutdown, nil
}
