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
	readingflows "rootstock/web-server/flows/reading"
	scoreflows "rootstock/web-server/flows/score"
	connecthandlers "rootstock/web-server/handlers/connect"
	httphandlers "rootstock/web-server/handlers/http"
	campaignops "rootstock/web-server/ops/campaign"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	mqttops "rootstock/web-server/ops/mqtt"
	orgops "rootstock/web-server/ops/org"
	readingops "rootstock/web-server/ops/reading"
	scoreops "rootstock/web-server/ops/score"
	"rootstock/web-server/proto/rootstock/v1/rootstockv1connect"
	"rootstock/web-server/repo/authorization"
	campaignrepo "rootstock/web-server/repo/campaign"
	identityrepo "rootstock/web-server/repo/identity"
	readingrepo "rootstock/web-server/repo/reading"
	scorerepo "rootstock/web-server/repo/score"
)

// NewRPCServer wires repos → ops → flows → Connect RPC handlers and returns
// an http.Handler, the MQTTFlows for subscription wiring, and a shutdown function.
func NewRPCServer(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, iRepo identityrepo.Repository, dOps *deviceops.Ops, crtOps *certops.Ops, mOps *mqttops.Ops) (http.Handler, *MQTTFlows, func(), error) {
	// JWT verification
	jwksURL := fmt.Sprintf("http://%s:%d/oauth/v2/keys", cfg.Identity.Zitadel.Host, cfg.Identity.Zitadel.Port)
	jwtVerifier, err := NewJWTVerifier(ctx, jwksURL, cfg.Identity.Zitadel.ExternalDomain, cfg.Identity.Zitadel.Issuer)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create jwt verifier: %w", err)
	}

	// Authorization (OPA)
	authzRepo := authorization.NewOPARepository()
	if err := authzRepo.Recompile(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("compile authorization policy: %w", err)
	}

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create otel interceptor: %w", err)
	}
	interceptors := connect.WithInterceptors(otelInterceptor, AuthorizationInterceptor(jwtVerifier, authzRepo), BinaryOnlyInterceptor())

	// Business repos
	cRepo := campaignrepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)
	sRepo := scorerepo.NewRepository(pool)

	// Ops
	cOps := campaignops.NewOps(cRepo)
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
	registerDeviceFlow := deviceflows.NewRegisterDeviceFlow(dOps, crtOps)
	getCACertFlow := deviceflows.NewGetCACertFlow(crtOps)
	enrollInCampaignFlow := deviceflows.NewEnrollInCampaignFlow(dOps, cOps, mOps)
	renewCertFlow := deviceflows.NewRenewCertFlow(dOps, crtOps)

	// Reading flows
	ingestReadingFlow := readingflows.NewIngestReadingFlow(cOps, rOps)
	exportDataFlow := readingflows.NewExportDataFlow(rOps)

	// Score flows
	getContributionFlow := scoreflows.NewGetContributionFlow(sOps)
	refreshScitizenScoreFlow := scoreflows.NewRefreshScitizenScoreFlow(dOps, rOps, sOps)

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
		exportDataFlow, cfg.Export.HMACSecret,
	)
	campaignPath, campaignH := rootstockv1connect.NewCampaignServiceHandler(campaignHandler, interceptors)

	orgHandler := connecthandlers.NewOrgServiceHandler(
		createOrgFlow, nestOrgFlow, defineRoleFlow, assignRoleFlow, inviteUserFlow,
	)
	orgPath, orgH := rootstockv1connect.NewOrgServiceHandler(orgHandler, interceptors)

	scoreHandler := connecthandlers.NewScoreServiceHandler(getContributionFlow)
	scorePath, scoreH := rootstockv1connect.NewScoreServiceHandler(scoreHandler, interceptors)

	deviceHandler := connecthandlers.NewDeviceServiceHandler(getDeviceFlow, revokeDeviceFlow, reinstateDeviceFlow, enrollInCampaignFlow)
	devicePath, deviceH := rootstockv1connect.NewDeviceServiceHandler(deviceHandler, interceptors)

	mux := http.NewServeMux()
	mux.Handle(healthPath, healthH)
	mux.Handle(campaignPath, campaignH)
	mux.Handle(orgPath, orgH)
	mux.Handle(scorePath, scoreH)
	mux.Handle(devicePath, deviceH)

	// Device enrollment (enrollment code auth, not JWT) + public CA cert
	enrollHandler := httphandlers.NewEnrollHandler(registerDeviceFlow, getCACertFlow)
	mux.HandleFunc("/enroll", enrollHandler.Enroll)
	mux.HandleFunc("/ca", enrollHandler.GetCACert)

	mqttFlows := &MQTTFlows{
		IngestReading:        ingestReadingFlow,
		RenewCert:            renewCertFlow,
		RefreshScitizenScore: refreshScitizenScoreFlow,
	}

	shutdown := func() {
		cRepo.Shutdown()
		rRepo.Shutdown()
		sRepo.Shutdown()
	}

	return mux, mqttFlows, shutdown, nil
}
