package connect

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	campaignflows "rootstock/web-server/flows/campaign"
	readingflows "rootstock/web-server/flows/reading"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// CampaignServiceHandler implements the CampaignService Connect RPC interface.
type CampaignServiceHandler struct {
	createCampaign    *campaignflows.CreateCampaignFlow
	publishCampaign   *campaignflows.PublishCampaignFlow
	browseCampaigns   *campaignflows.BrowseCampaignsFlow
	campaignDashboard *campaignflows.DashboardFlow
	exportData        *readingflows.ExportDataFlow
	hmacSecret        string
}

// NewCampaignServiceHandler creates the handler with all required flows.
func NewCampaignServiceHandler(
	createCampaign *campaignflows.CreateCampaignFlow,
	publishCampaign *campaignflows.PublishCampaignFlow,
	browseCampaigns *campaignflows.BrowseCampaignsFlow,
	campaignDashboard *campaignflows.DashboardFlow,
	exportData *readingflows.ExportDataFlow,
	hmacSecret string,
) *CampaignServiceHandler {
	return &CampaignServiceHandler{
		createCampaign:    createCampaign,
		publishCampaign:   publishCampaign,
		browseCampaigns:   browseCampaigns,
		campaignDashboard: campaignDashboard,
		exportData:        exportData,
		hmacSecret:        hmacSecret,
	}
}

func (h *CampaignServiceHandler) CreateCampaign(
	ctx context.Context,
	req *connect.Request[rootstockv1.CreateCampaignRequest],
) (*connect.Response[rootstockv1.CreateCampaignResponse], error) {
	msg := req.Msg

	input := campaignflows.CreateCampaignInput{
		OrgID:       msg.GetOrgId(),
		CreatedBy:   msg.GetCreatedBy(),
		WindowStart: parseOptionalTime(msg.WindowStart),
		WindowEnd:   parseOptionalTime(msg.WindowEnd),
	}

	for _, p := range msg.GetParameters() {
		pi := campaignflows.ParameterInput{
			Name:     p.GetName(),
			Unit:     p.GetUnit(),
			MinRange: p.MinRange,
			MaxRange: p.MaxRange,
		}
		if p.Precision != nil {
			v := int(p.GetPrecision())
			pi.Precision = &v
		}
		input.Parameters = append(input.Parameters, pi)
	}
	for _, r := range msg.GetRegions() {
		input.Regions = append(input.Regions, campaignflows.RegionInput{GeoJSON: r.GetGeoJson()})
	}
	for _, e := range msg.GetEligibility() {
		input.Eligibility = append(input.Eligibility, campaignflows.EligibilityInput{
			DeviceClass:     e.GetDeviceClass(),
			Tier:            int(e.GetTier()),
			RequiredSensors: e.GetRequiredSensors(),
			FirmwareMin:     e.GetFirmwareMin(),
		})
	}

	result, err := h.createCampaign.Run(ctx, input)
	if err != nil {
		slog.ErrorContext(ctx, "create campaign failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&rootstockv1.CreateCampaignResponse{
		Campaign: campaignToProto(result),
	}), nil
}

func (h *CampaignServiceHandler) PublishCampaign(
	ctx context.Context,
	req *connect.Request[rootstockv1.PublishCampaignRequest],
) (*connect.Response[rootstockv1.PublishCampaignResponse], error) {
	if err := h.publishCampaign.Run(ctx, req.Msg.GetCampaignId()); err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.PublishCampaignResponse{}), nil
}

func (h *CampaignServiceHandler) ListCampaigns(
	ctx context.Context,
	req *connect.Request[rootstockv1.ListCampaignsRequest],
) (*connect.Response[rootstockv1.ListCampaignsResponse], error) {
	msg := req.Msg

	input := campaignflows.BrowseCampaignsInput{
		Status:    msg.GetStatus(),
		OrgID:     msg.GetOrgId(),
		Longitude: msg.Longitude,
		Latitude:  msg.Latitude,
		RadiusKm:  msg.RadiusKm,
	}

	campaigns, err := h.browseCampaigns.Run(ctx, input)
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.CampaignProto, len(campaigns))
	for i := range campaigns {
		protos[i] = campaignToProto(&campaigns[i])
	}

	return connect.NewResponse(&rootstockv1.ListCampaignsResponse{
		Campaigns: protos,
	}), nil
}

func (h *CampaignServiceHandler) GetCampaignDashboard(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetCampaignDashboardRequest],
) (*connect.Response[rootstockv1.GetCampaignDashboardResponse], error) {
	dashboard, err := h.campaignDashboard.Run(ctx, req.Msg.GetCampaignId())
	if err != nil {
		return nil, err
	}

	resp := &rootstockv1.GetCampaignDashboardResponse{
		CampaignId:      dashboard.CampaignID,
		AcceptedCount:   int32(dashboard.AcceptedCount),
		QuarantineCount: int32(dashboard.QuarantineCount),
	}

	for _, pq := range dashboard.ParameterQuality {
		resp.ParameterQuality = append(resp.ParameterQuality, &rootstockv1.ParameterQualityProto{
			ParameterName:    pq.ParameterName,
			AcceptedCount:    int32(pq.AcceptedCount),
			QuarantinedCount: int32(pq.QuarantinedCount),
		})
	}

	for _, db := range dashboard.DeviceBreakdown {
		entry := &rootstockv1.DeviceBreakdownProto{
			PseudoDeviceId: db.PseudoDeviceID,
			DeviceClass:    db.DeviceClass,
			AcceptanceRate: db.AcceptanceRate,
			ReadingCount:   int32(db.ReadingCount),
		}
		if db.LastSeen != nil {
			entry.LastSeen = *db.LastSeen
		}
		resp.DeviceBreakdown = append(resp.DeviceBreakdown, entry)
	}

	resp.EnrollmentFunnel = &rootstockv1.EnrollmentFunnelProto{
		Enrolled:     int32(dashboard.EnrollmentFunnel.Enrolled),
		Active:       int32(dashboard.EnrollmentFunnel.Active),
		Contributing: int32(dashboard.EnrollmentFunnel.Contributing),
	}

	for _, tc := range dashboard.TemporalCoverage {
		resp.TemporalCoverage = append(resp.TemporalCoverage, &rootstockv1.TemporalBucketProto{
			Bucket: tc.Bucket,
			Count:  int32(tc.Count),
		})
	}

	return connect.NewResponse(resp), nil
}

func (h *CampaignServiceHandler) ExportCampaignData(
	ctx context.Context,
	req *connect.Request[rootstockv1.ExportCampaignDataRequest],
) (*connect.Response[rootstockv1.ExportCampaignDataResponse], error) {
	msg := req.Msg

	result, err := h.exportData.Run(ctx, readingflows.ExportDataInput{
		CampaignID: msg.GetCampaignId(),
		Secret:     h.hmacSecret,
		Limit:      int(msg.GetLimit()),
		Offset:     int(msg.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	readings := make([]*rootstockv1.ExportedReadingProto, len(result.Readings))
	for i, r := range result.Readings {
		readings[i] = &rootstockv1.ExportedReadingProto{
			PseudoDeviceId:  r.PseudoDeviceID,
			CampaignId:      r.CampaignID,
			Values:          r.Values,
			Timestamp:       r.Timestamp.Format(time.RFC3339),
			FirmwareVersion: r.FirmwareVersion,
			IngestedAt:      r.IngestedAt.Format(time.RFC3339),
			Status:          r.Status,
		}
		if r.Geolocation != nil {
			readings[i].Geolocation = r.Geolocation
		}
	}

	return connect.NewResponse(&rootstockv1.ExportCampaignDataResponse{
		Readings: readings,
	}), nil
}

func campaignToProto(c *campaignflows.Campaign) *rootstockv1.CampaignProto {
	proto := &rootstockv1.CampaignProto{
		Id:        c.ID,
		OrgId:     c.OrgID,
		Status:    c.Status,
		CreatedBy: c.CreatedBy,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
	if c.WindowStart != nil {
		s := c.WindowStart.Format(time.RFC3339)
		proto.WindowStart = &s
	}
	if c.WindowEnd != nil {
		s := c.WindowEnd.Format(time.RFC3339)
		proto.WindowEnd = &s
	}
	return proto
}

func parseOptionalTime(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return &t
}
