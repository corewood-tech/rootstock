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
	scitizenflows "rootstock/web-server/flows/scitizen"
	scoreflows "rootstock/web-server/flows/score"
	securityflows "rootstock/web-server/flows/security"
	userflows "rootstock/web-server/flows/user"
	connecthandlers "rootstock/web-server/handlers/connect"
	httphandlers "rootstock/web-server/handlers/http"
	campaignops "rootstock/web-server/ops/campaign"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	enrollmentops "rootstock/web-server/ops/enrollment"
	graphops "rootstock/web-server/ops/graph"
	mqttops "rootstock/web-server/ops/mqtt"
	notificationops "rootstock/web-server/ops/notification"
	orgops "rootstock/web-server/ops/org"
	readingops "rootstock/web-server/ops/reading"
	scitizenops "rootstock/web-server/ops/scitizen"
	scoreops "rootstock/web-server/ops/score"
	userops "rootstock/web-server/ops/user"
	"rootstock/web-server/proto/rootstock/v1/rootstockv1connect"
	"rootstock/web-server/repo/authorization"
	campaignrepo "rootstock/web-server/repo/campaign"
	enrollmentrepo "rootstock/web-server/repo/enrollment"
	graphrepo "rootstock/web-server/repo/graph"
	identityrepo "rootstock/web-server/repo/identity"
	notificationrepo "rootstock/web-server/repo/notification"
	readingrepo "rootstock/web-server/repo/reading"
	scitizenrepo "rootstock/web-server/repo/scitizen"
	scorerepo "rootstock/web-server/repo/score"
	sessionrepo "rootstock/web-server/repo/session"
	userrepo "rootstock/web-server/repo/user"
)

// NewRPCServer wires repos → ops → flows → Connect RPC handlers and returns
// an http.Handler, the MQTTFlows for subscription wiring, and a shutdown function.
func NewRPCServer(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, iRepo identityrepo.Repository, dOps *deviceops.Ops, crtOps *certops.Ops, mOps *mqttops.Ops) (http.Handler, *MQTTFlows, func(), error) {
	// Session repo (Zitadel Session API)
	sessRepo, err := sessionrepo.NewRepository(ctx, cfg.Identity.Zitadel)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create session repo: %w", err)
	}

	// Authorization (OPA)
	authzRepo := authorization.NewOPARepository()
	if err := authzRepo.Recompile(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("compile authorization policy: %w", err)
	}

	// Business repos
	cRepo := campaignrepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)
	sRepo := scorerepo.NewRepository(pool)
	uRepo := userrepo.NewRepository(pool)
	scRepo := scitizenrepo.NewRepository(pool)
	eRepo := enrollmentrepo.NewRepository(pool)

	// Notification repo (SMTP)
	nRepo := notificationrepo.NewRepository(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From)

	// Graph repo (Dgraph)
	gRepo, err := graphrepo.NewDgraphRepository(cfg.Database.Dgraph.AlphaAddr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create graph repository: %w", err)
	}

	// Ops
	cOps := campaignops.NewOps(cRepo)
	rOps := readingops.NewOps(rRepo)
	oOps := orgops.NewOps(iRepo)
	sOps := scoreops.NewOps(sRepo)
	uOps := userops.NewOps(uRepo, sessRepo)
	nOps := notificationops.NewOps(nRepo)
	gOps := graphops.NewOps(gRepo)
	scOps := scitizenops.NewOps(scRepo)
	eOps := enrollmentops.NewOps(eRepo)

	// Interceptors (session-based auth)
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create otel interceptor: %w", err)
	}
	interceptors := connect.WithInterceptors(otelInterceptor, AuthorizationInterceptor(uOps, authzRepo), BinaryOnlyInterceptor())

	// Campaign flows
	createCampaignFlow := campaignflows.NewCreateCampaignFlow(cOps, gOps)
	publishCampaignFlow := campaignflows.NewPublishCampaignFlow(cOps, gOps)
	browseCampaignsFlow := campaignflows.NewBrowseCampaignsFlow(cOps)
	campaignDashboardFlow := campaignflows.NewDashboardFlow(rOps)

	// Device flows
	getDeviceFlow := deviceflows.NewGetDeviceFlow(dOps)
	revokeDeviceFlow := deviceflows.NewRevokeDeviceFlow(dOps)
	reinstateDeviceFlow := deviceflows.NewReinstateDeviceFlow(dOps)
	registerDeviceFlow := deviceflows.NewRegisterDeviceFlow(dOps, crtOps)
	getCACertFlow := deviceflows.NewGetCACertFlow(crtOps)
	enrollInCampaignFlow := deviceflows.NewEnrollInCampaignFlow(dOps, cOps, mOps, gOps)
	renewCertFlow := deviceflows.NewRenewCertFlow(dOps, crtOps)

	// Reading flows
	ingestReadingFlow := readingflows.NewIngestReadingFlow(cOps, rOps, gOps)
	exportDataFlow := readingflows.NewExportDataFlow(rOps)

	// Score flows
	getContributionFlow := scoreflows.NewGetContributionFlow(sOps)
	refreshScitizenScoreFlow := scoreflows.NewRefreshScitizenScoreFlow(dOps, rOps, sOps)

	// Security flows
	securityResponseFlow := securityflows.NewSecurityResponseFlow(dOps, rOps, nOps)

	// User flows
	registerUserFlow := userflows.NewRegisterUserFlow(uOps)
	getUserFlow := userflows.NewGetUserFlow(uOps)
	loginFlow := userflows.NewLoginFlow(uOps)
	logoutFlow := userflows.NewLogoutFlow(uOps)
	registerResearcherFlow := userflows.NewRegisterResearcherFlow(oOps, uOps, nOps, cfg.Server.PublicBaseURL)
	verifyEmailFlow := userflows.NewVerifyEmailFlow(oOps)
	updateUserTypeFlow := userflows.NewUpdateUserTypeFlow(uOps)

	// Scitizen flows
	scitizenRegisterFlow := scitizenflows.NewScitzenRegistrationFlow(scOps)
	scitizenDashboardFlow := scitizenflows.NewScitzenDashboardFlow(scOps)
	scitizenBrowseCampaignsFlow := scitizenflows.NewBrowseCampaignsFlow(scOps)
	scitizenCampaignDetailFlow := scitizenflows.NewCampaignDetailFlow(scOps)
	scitizenCampaignSearchFlow := scitizenflows.NewCampaignSearchFlow(scOps)
	scitizenEnrollDeviceFlow := scitizenflows.NewEnrollDeviceCampaignFlow(scOps, eOps)
	scitizenWithdrawFlow := scitizenflows.NewWithdrawEnrollmentFlow(eOps)
	scitizenDeviceFlow := scitizenflows.NewDeviceManagementFlow(scOps)
	scitizenOnboardingFlow := scitizenflows.NewOnboardingFlow(scOps)
	scitizenNotificationFlow := scitizenflows.NewNotificationFlow(scOps)
	scitizenProgressFlow := scitizenflows.NewCampaignProgressFlow(scOps)

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

	userHandler := connecthandlers.NewUserServiceHandler(registerUserFlow, getUserFlow, loginFlow, logoutFlow, registerResearcherFlow, verifyEmailFlow, updateUserTypeFlow)
	userPath, userH := rootstockv1connect.NewUserServiceHandler(userHandler, interceptors)

	scitizenHandler := connecthandlers.NewScitizenServiceHandler(
		getUserFlow, scitizenRegisterFlow, scitizenDashboardFlow,
		scitizenBrowseCampaignsFlow, scitizenCampaignDetailFlow, scitizenCampaignSearchFlow,
		scitizenEnrollDeviceFlow, scitizenWithdrawFlow, scitizenDeviceFlow,
		scitizenOnboardingFlow, scitizenNotificationFlow, scitizenProgressFlow,
	)
	scitizenPath, scitizenH := rootstockv1connect.NewScitizenServiceHandler(scitizenHandler, interceptors)

	adminHandler := connecthandlers.NewAdminServiceHandler(securityResponseFlow)
	adminPath, adminH := rootstockv1connect.NewAdminServiceHandler(adminHandler, interceptors)

	mux := http.NewServeMux()
	mux.Handle(healthPath, healthH)
	mux.Handle(campaignPath, campaignH)
	mux.Handle(orgPath, orgH)
	mux.Handle(scorePath, scoreH)
	mux.Handle(devicePath, deviceH)
	mux.Handle(userPath, userH)
	mux.Handle(scitizenPath, scitizenH)
	mux.Handle(adminPath, adminH)

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
		uRepo.Shutdown()
		nRepo.Shutdown()
		gRepo.Shutdown()
		scRepo.Shutdown()
		eRepo.Shutdown()
		sessRepo.Shutdown()
	}

	return mux, mqttFlows, shutdown, nil
}
