package connect

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"rootstock/web-server/auth"
	scitizenflows "rootstock/web-server/flows/scitizen"
	userflows "rootstock/web-server/flows/user"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// ScitizenServiceHandler implements the ScitizenService Connect RPC interface.
// Graph node: 0x2d (ScitzenServiceHandler)
type ScitizenServiceHandler struct {
	getUser            *userflows.GetUserFlow
	register           *scitizenflows.ScitzenRegistrationFlow
	dashboard          *scitizenflows.ScitzenDashboardFlow
	browseCampaigns    *scitizenflows.BrowseCampaignsFlow
	campaignDetail     *scitizenflows.CampaignDetailFlow
	campaignSearch     *scitizenflows.CampaignSearchFlow
	enrollDevice       *scitizenflows.EnrollDeviceCampaignFlow
	withdrawEnrollment *scitizenflows.WithdrawEnrollmentFlow
	deviceManagement   *scitizenflows.DeviceManagementFlow
	onboarding         *scitizenflows.OnboardingFlow
	notifications      *scitizenflows.NotificationFlow
	campaignProgress   *scitizenflows.CampaignProgressFlow
}

// NewScitizenServiceHandler creates the handler with all required flows.
func NewScitizenServiceHandler(
	getUser *userflows.GetUserFlow,
	register *scitizenflows.ScitzenRegistrationFlow,
	dashboard *scitizenflows.ScitzenDashboardFlow,
	browseCampaigns *scitizenflows.BrowseCampaignsFlow,
	campaignDetail *scitizenflows.CampaignDetailFlow,
	campaignSearch *scitizenflows.CampaignSearchFlow,
	enrollDevice *scitizenflows.EnrollDeviceCampaignFlow,
	withdrawEnrollment *scitizenflows.WithdrawEnrollmentFlow,
	deviceManagement *scitizenflows.DeviceManagementFlow,
	onboarding *scitizenflows.OnboardingFlow,
	notifications *scitizenflows.NotificationFlow,
	campaignProgress *scitizenflows.CampaignProgressFlow,
) *ScitizenServiceHandler {
	return &ScitizenServiceHandler{
		getUser:            getUser,
		register:           register,
		dashboard:          dashboard,
		browseCampaigns:    browseCampaigns,
		campaignDetail:     campaignDetail,
		campaignSearch:     campaignSearch,
		enrollDevice:       enrollDevice,
		withdrawEnrollment: withdrawEnrollment,
		deviceManagement:   deviceManagement,
		onboarding:         onboarding,
		notifications:      notifications,
		campaignProgress:   campaignProgress,
	}
}

// resolveUserID extracts the IdP user ID from context and resolves the app user ID.
func (h *ScitizenServiceHandler) resolveUserID(ctx context.Context) (string, error) {
	idpID, ok := auth.SubjectFromContext(ctx)
	if !ok || idpID == "" {
		return "", connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no authenticated subject"))
	}
	user, err := h.getUser.Run(ctx, idpID)
	if err != nil {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("resolve user: %w", err))
	}
	if user == nil {
		return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}
	return user.ID, nil
}

func (h *ScitizenServiceHandler) RegisterScitizen(
	ctx context.Context,
	req *connect.Request[rootstockv1.RegisterScitizenRequest],
) (*connect.Response[rootstockv1.RegisterScitizenResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	_, err = h.register.Run(ctx, scitizenflows.RegisterInput{
		UserID:     userID,
		TOSVersion: "1.0",
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("register scitizen: %w", err))
	}

	return connect.NewResponse(&rootstockv1.RegisterScitizenResponse{
		UserId:                userID,
		EmailVerificationSent: true,
	}), nil
}

func (h *ScitizenServiceHandler) GetDashboard(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetDashboardRequest],
) (*connect.Response[rootstockv1.GetDashboardResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.dashboard.Run(ctx, userID)
	if err != nil {
		return nil, err
	}

	badges := make([]*rootstockv1.BadgeProto, len(result.Badges))
	for i, b := range result.Badges {
		badges[i] = &rootstockv1.BadgeProto{
			Id:        b.ID,
			BadgeType: b.BadgeType,
			AwardedAt: b.AwardedAt.Format(time.RFC3339),
		}
	}

	enrollments := make([]*rootstockv1.EnrollmentProto, len(result.Enrollments))
	for i, e := range result.Enrollments {
		enrollments[i] = &rootstockv1.EnrollmentProto{
			Id:         e.ID,
			DeviceId:   e.DeviceID,
			CampaignId: e.CampaignID,
			Status:     e.Status,
			EnrolledAt: e.EnrolledAt.Format(time.RFC3339),
		}
	}

	return connect.NewResponse(&rootstockv1.GetDashboardResponse{
		ActiveEnrollments: int32(result.ActiveEnrollments),
		TotalReadings:     int32(result.TotalReadings),
		AcceptedReadings:  int32(result.AcceptedReadings),
		ContributionScore: result.ContributionScore,
		Badges:            badges,
		Enrollments:       enrollments,
	}), nil
}

func (h *ScitizenServiceHandler) BrowsePublishedCampaigns(
	ctx context.Context,
	req *connect.Request[rootstockv1.BrowsePublishedCampaignsRequest],
) (*connect.Response[rootstockv1.BrowsePublishedCampaignsResponse], error) {
	msg := req.Msg

	input := scitizenflows.BrowseInput{
		Longitude: msg.Longitude,
		Latitude:  msg.Latitude,
		RadiusKm:  msg.RadiusKm,
		Limit:     int(msg.GetLimit()),
		Offset:    int(msg.GetOffset()),
	}
	if msg.SensorType != nil {
		input.SensorType = msg.SensorType
	}

	campaigns, total, err := h.browseCampaigns.Run(ctx, input)
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.CampaignSummaryProto, len(campaigns))
	for i, c := range campaigns {
		protos[i] = campaignSummaryToProto(&c)
	}

	return connect.NewResponse(&rootstockv1.BrowsePublishedCampaignsResponse{
		Campaigns: protos,
		Total:     int32(total),
	}), nil
}

func (h *ScitizenServiceHandler) GetCampaignDetail(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetCampaignDetailRequest],
) (*connect.Response[rootstockv1.GetCampaignDetailResponse], error) {
	result, err := h.campaignDetail.Run(ctx, req.Msg.GetCampaignId())
	if err != nil {
		return nil, err
	}

	params := make([]*rootstockv1.ParameterProto, len(result.Parameters))
	for i, p := range result.Parameters {
		params[i] = &rootstockv1.ParameterProto{
			Name:      p.Name,
			Unit:      p.Unit,
			MinRange:  p.MinRange,
			MaxRange:  p.MaxRange,
		}
		if p.Precision != nil {
			v := int32(*p.Precision)
			params[i].Precision = &v
		}
	}

	regions := make([]*rootstockv1.RegionProto, len(result.Regions))
	for i, r := range result.Regions {
		regions[i] = &rootstockv1.RegionProto{GeoJson: r.GeoJSON}
	}

	elig := make([]*rootstockv1.EligibilityProto, len(result.Eligibility))
	for i, e := range result.Eligibility {
		elig[i] = &rootstockv1.EligibilityProto{
			DeviceClass:     e.DeviceClass,
			Tier:            int32(e.Tier),
			RequiredSensors: e.RequiredSensors,
			FirmwareMin:     e.FirmwareMin,
		}
	}

	resp := &rootstockv1.GetCampaignDetailResponse{
		CampaignId:      result.CampaignID,
		Status:          result.Status,
		Parameters:      params,
		Regions:         regions,
		Eligibility:     elig,
		EnrollmentCount: int32(result.EnrollmentCount),
		ProgressPercent: result.ProgressPercent,
	}
	if result.WindowStart != nil {
		s := result.WindowStart.Format(time.RFC3339)
		resp.WindowStart = &s
	}
	if result.WindowEnd != nil {
		s := result.WindowEnd.Format(time.RFC3339)
		resp.WindowEnd = &s
	}

	return connect.NewResponse(resp), nil
}

func (h *ScitizenServiceHandler) SearchCampaigns(
	ctx context.Context,
	req *connect.Request[rootstockv1.SearchCampaignsRequest],
) (*connect.Response[rootstockv1.SearchCampaignsResponse], error) {
	campaigns, total, err := h.campaignSearch.Run(ctx, scitizenflows.SearchInput{
		Query:  req.Msg.GetQuery(),
		Limit:  int(req.Msg.GetLimit()),
		Offset: int(req.Msg.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.CampaignSummaryProto, len(campaigns))
	for i, c := range campaigns {
		protos[i] = campaignSummaryToProto(&c)
	}

	return connect.NewResponse(&rootstockv1.SearchCampaignsResponse{
		Campaigns: protos,
		Total:     int32(total),
	}), nil
}

func (h *ScitizenServiceHandler) EnrollDevice(
	ctx context.Context,
	req *connect.Request[rootstockv1.EnrollDeviceRequest],
) (*connect.Response[rootstockv1.EnrollDeviceResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	msg := req.Msg
	consent := msg.GetConsent()

	result, err := h.enrollDevice.Run(ctx, scitizenflows.EnrollDeviceInput{
		ScitizenID:     userID,
		DeviceID:       msg.GetDeviceId(),
		CampaignID:     msg.GetCampaignId(),
		ConsentVersion: consent.GetVersion(),
		ConsentScope:   consent.GetScope(),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.EnrollDeviceResponse{
		Enrolled:     result.Enrolled,
		Reason:       result.Reason,
		EnrollmentId: result.EnrollmentID,
	}), nil
}

func (h *ScitizenServiceHandler) WithdrawEnrollment(
	ctx context.Context,
	req *connect.Request[rootstockv1.WithdrawEnrollmentRequest],
) (*connect.Response[rootstockv1.WithdrawEnrollmentResponse], error) {
	if err := h.withdrawEnrollment.Run(ctx, req.Msg.GetEnrollmentId()); err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.WithdrawEnrollmentResponse{}), nil
}

func (h *ScitizenServiceHandler) GetDevices(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetDevicesRequest],
) (*connect.Response[rootstockv1.GetDevicesResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	devices, err := h.deviceManagement.RunList(ctx, userID)
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.DeviceSummaryProto, len(devices))
	for i, d := range devices {
		protos[i] = &rootstockv1.DeviceSummaryProto{
			Id:                d.ID,
			Status:            d.Status,
			Class:             d.Class,
			FirmwareVersion:   d.FirmwareVersion,
			Tier:              int32(d.Tier),
			Sensors:           d.Sensors,
			ActiveEnrollments: int32(d.ActiveEnrollments),
		}
		if d.LastSeen != nil {
			s := d.LastSeen.Format(time.RFC3339)
			protos[i].LastSeen = &s
		}
	}

	return connect.NewResponse(&rootstockv1.GetDevicesResponse{
		Devices: protos,
	}), nil
}

func (h *ScitizenServiceHandler) GetDeviceDetail(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetDeviceDetailRequest],
) (*connect.Response[rootstockv1.GetDeviceDetailResponse], error) {
	result, err := h.deviceManagement.RunDetail(ctx, req.Msg.GetDeviceId())
	if err != nil {
		return nil, err
	}

	enrollments := make([]*rootstockv1.EnrollmentProto, len(result.Enrollments))
	for i, e := range result.Enrollments {
		enrollments[i] = &rootstockv1.EnrollmentProto{
			Id:         e.ID,
			DeviceId:   e.DeviceID,
			CampaignId: e.CampaignID,
			Status:     e.Status,
			EnrolledAt: e.EnrolledAt.Format(time.RFC3339),
		}
	}

	connHistory := make([]*rootstockv1.ConnectionEventProto, len(result.ConnectionHistory))
	for i, c := range result.ConnectionHistory {
		connHistory[i] = &rootstockv1.ConnectionEventProto{
			EventType: c.EventType,
			Timestamp: c.Timestamp.Format(time.RFC3339),
			Reason:    c.Reason,
		}
	}

	device := &rootstockv1.DeviceProto{
		Id:              result.ID,
		OwnerId:         result.OwnerID,
		Status:          result.Status,
		Class:           result.Class,
		FirmwareVersion: result.FirmwareVersion,
		Tier:            int32(result.Tier),
		Sensors:         result.Sensors,
		CertSerial:      result.CertSerial,
		CreatedAt:       result.CreatedAt.Format(time.RFC3339),
	}

	return connect.NewResponse(&rootstockv1.GetDeviceDetailResponse{
		Device:            device,
		Enrollments:       enrollments,
		ConnectionHistory: connHistory,
	}), nil
}

func (h *ScitizenServiceHandler) GetNotifications(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetNotificationsRequest],
) (*connect.Response[rootstockv1.GetNotificationsResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	msg := req.Msg
	result, err := h.notifications.Run(ctx, scitizenflows.GetNotificationsInput{
		UserID:     userID,
		TypeFilter: msg.TypeFilter,
		Limit:      int(msg.GetLimit()),
		Offset:     int(msg.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.NotificationProto, len(result.Notifications))
	for i, n := range result.Notifications {
		protos[i] = &rootstockv1.NotificationProto{
			Id:           n.ID,
			Type:         n.Type,
			Message:      n.Message,
			Read:         n.Read,
			ResourceLink: n.ResourceLink,
			CreatedAt:    n.CreatedAt.Format(time.RFC3339),
		}
	}

	return connect.NewResponse(&rootstockv1.GetNotificationsResponse{
		Notifications: protos,
		UnreadCount:   int32(result.UnreadCount),
		Total:         int32(result.Total),
	}), nil
}

func (h *ScitizenServiceHandler) GetContributions(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetContributionsRequest],
) (*connect.Response[rootstockv1.GetContributionsResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.campaignProgress.Run(ctx, userID)
	if err != nil {
		return nil, err
	}

	histories := make([]*rootstockv1.ReadingHistoryProto, len(result.Histories))
	for i, h := range result.Histories {
		histories[i] = &rootstockv1.ReadingHistoryProto{
			DeviceId:      h.DeviceID,
			CampaignId:    h.CampaignID,
			TotalReadings: int32(h.Total),
			Accepted:      int32(h.Accepted),
			Rejected:      int32(h.Rejected),
		}
	}

	badges := make([]*rootstockv1.BadgeProto, len(result.Badges))
	for i, b := range result.Badges {
		badges[i] = &rootstockv1.BadgeProto{
			Id:        b.ID,
			BadgeType: b.BadgeType,
			AwardedAt: b.AwardedAt.Format(time.RFC3339),
		}
	}

	return connect.NewResponse(&rootstockv1.GetContributionsResponse{
		Histories:         histories,
		ContributionScore: result.ContributionScore,
		Badges:            badges,
	}), nil
}

func (h *ScitizenServiceHandler) GetOnboardingState(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetOnboardingStateRequest],
) (*connect.Response[rootstockv1.GetOnboardingStateResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	state, err := h.onboarding.Run(ctx, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.GetOnboardingStateResponse{
		State: &rootstockv1.OnboardingStateProto{
			DeviceRegistered:     state.DeviceRegistered,
			CampaignEnrolled:     state.CampaignEnrolled,
			FirstReadingSubmitted: state.FirstReadingSubmitted,
			TosAccepted:          state.TOSAccepted,
		},
	}), nil
}

func campaignSummaryToProto(c *scitizenflows.CampaignSummary) *rootstockv1.CampaignSummaryProto {
	p := &rootstockv1.CampaignSummaryProto{
		Id:              c.ID,
		Status:          c.Status,
		EnrollmentCount: int32(c.EnrollmentCount),
		RequiredSensors: c.RequiredSensors,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
	}
	if c.WindowStart != nil {
		s := c.WindowStart.Format(time.RFC3339)
		p.WindowStart = &s
	}
	if c.WindowEnd != nil {
		s := c.WindowEnd.Format(time.RFC3339)
		p.WindowEnd = &s
	}
	return p
}
