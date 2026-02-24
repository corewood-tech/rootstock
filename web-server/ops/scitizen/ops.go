package scitizen

import (
	"context"

	scitizenrepo "rootstock/web-server/repo/scitizen"
)

// Ops holds scitizen operations. Each method is one op.
// Graph node: 0x2a (ScitizenOps)
type Ops struct {
	repo scitizenrepo.Repository
}

// NewOps creates scitizen ops backed by the given repository.
func NewOps(repo scitizenrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// CreateProfile creates a scitizen profile with ToS acceptance.
func (o *Ops) CreateProfile(ctx context.Context, input CreateProfileInput) (*Profile, error) {
	result, err := o.repo.CreateProfile(ctx, scitizenrepo.CreateProfileInput{
		UserID:     input.UserID,
		TOSVersion: input.TOSVersion,
	})
	if err != nil {
		return nil, err
	}
	return fromRepoProfile(result), nil
}

// GetProfile returns the scitizen profile.
func (o *Ops) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	result, err := o.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return fromRepoProfile(result), nil
}

// UpdateOnboarding updates onboarding state flags.
func (o *Ops) UpdateOnboarding(ctx context.Context, input UpdateOnboardingInput) error {
	return o.repo.UpdateOnboarding(ctx, scitizenrepo.UpdateOnboardingInput{
		UserID:           input.UserID,
		DeviceRegistered: input.DeviceRegistered,
		CampaignEnrolled: input.CampaignEnrolled,
		FirstReading:     input.FirstReading,
	})
}

// GetDashboard returns aggregated dashboard data.
func (o *Ops) GetDashboard(ctx context.Context, userID string) (*Dashboard, error) {
	result, err := o.repo.GetDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}
	return fromRepoDashboard(result), nil
}

// GetContributions returns reading history per device per campaign.
func (o *Ops) GetContributions(ctx context.Context, userID string) ([]ReadingHistory, error) {
	results, err := o.repo.GetContributions(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]ReadingHistory, len(results))
	for i, r := range results {
		out[i] = ReadingHistory{
			DeviceID:   r.DeviceID,
			CampaignID: r.CampaignID,
			Total:      r.Total,
			Accepted:   r.Accepted,
			Rejected:   r.Rejected,
		}
	}
	return out, nil
}

// BrowseCampaigns returns published campaigns with filtering.
func (o *Ops) BrowseCampaigns(ctx context.Context, input BrowseInput) ([]CampaignSummary, int, error) {
	results, total, err := o.repo.BrowseCampaigns(ctx, scitizenrepo.BrowseInput{
		Longitude:  input.Longitude,
		Latitude:   input.Latitude,
		RadiusKm:   input.RadiusKm,
		SensorType: input.SensorType,
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]CampaignSummary, len(results))
	for i, r := range results {
		out[i] = fromRepoCampaignSummary(&r)
	}
	return out, total, nil
}

// GetCampaignDetail returns full campaign detail for enrollment decision.
func (o *Ops) GetCampaignDetail(ctx context.Context, campaignID string) (*CampaignDetail, error) {
	result, err := o.repo.GetCampaignDetail(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return fromRepoCampaignDetail(result), nil
}

// SearchCampaigns performs full-text search across published campaigns.
func (o *Ops) SearchCampaigns(ctx context.Context, input SearchInput) ([]CampaignSummary, int, error) {
	results, total, err := o.repo.SearchCampaigns(ctx, scitizenrepo.SearchInput{
		Query:  input.Query,
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]CampaignSummary, len(results))
	for i, r := range results {
		out[i] = fromRepoCampaignSummary(&r)
	}
	return out, total, nil
}

// GetDevices returns all devices owned by the scitizen.
func (o *Ops) GetDevices(ctx context.Context, ownerID string) ([]DeviceSummary, error) {
	results, err := o.repo.GetDevices(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]DeviceSummary, len(results))
	for i, r := range results {
		out[i] = DeviceSummary{
			ID:                r.ID,
			Status:            r.Status,
			Class:             r.Class,
			FirmwareVersion:   r.FirmwareVersion,
			Tier:              r.Tier,
			Sensors:           r.Sensors,
			ActiveEnrollments: r.ActiveEnrollments,
			LastSeen:          r.LastSeen,
		}
	}
	return out, nil
}

// GetDeviceDetail returns full device info with enrollments and connection history.
func (o *Ops) GetDeviceDetail(ctx context.Context, deviceID string) (*DeviceDetail, error) {
	result, err := o.repo.GetDeviceDetail(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	return fromRepoDeviceDetail(result), nil
}

// GetNotifications returns notifications for the scitizen.
func (o *Ops) GetNotifications(ctx context.Context, input GetNotificationsInput) ([]Notification, int, int, error) {
	results, unreadCount, total, err := o.repo.GetNotifications(ctx, scitizenrepo.GetNotificationsInput{
		UserID:     input.UserID,
		TypeFilter: input.TypeFilter,
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, 0, 0, err
	}
	out := make([]Notification, len(results))
	for i, r := range results {
		out[i] = Notification{
			ID:           r.ID,
			UserID:       r.UserID,
			Type:         r.Type,
			Message:      r.Message,
			Read:         r.Read,
			ResourceLink: r.ResourceLink,
			CreatedAt:    r.CreatedAt,
		}
	}
	return out, unreadCount, total, nil
}

// --- converters ---

func fromRepoProfile(r *scitizenrepo.Profile) *Profile {
	return &Profile{
		UserID:           r.UserID,
		TOSAccepted:      r.TOSAccepted,
		TOSVersion:       r.TOSVersion,
		TOSAcceptedAt:    r.TOSAcceptedAt,
		DeviceRegistered: r.DeviceRegistered,
		CampaignEnrolled: r.CampaignEnrolled,
		FirstReading:     r.FirstReading,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func fromRepoDashboard(r *scitizenrepo.Dashboard) *Dashboard {
	badges := make([]Badge, len(r.Badges))
	for i, b := range r.Badges {
		badges[i] = Badge{ID: b.ID, BadgeType: b.BadgeType, AwardedAt: b.AwardedAt}
	}
	enrollments := make([]Enrollment, len(r.Enrollments))
	for i, e := range r.Enrollments {
		enrollments[i] = Enrollment{
			ID: e.ID, DeviceID: e.DeviceID, CampaignID: e.CampaignID,
			Status: e.Status, EnrolledAt: e.EnrolledAt,
		}
	}
	return &Dashboard{
		ActiveEnrollments: r.ActiveEnrollments,
		TotalReadings:     r.TotalReadings,
		AcceptedReadings:  r.AcceptedReadings,
		ContributionScore: r.ContributionScore,
		Badges:            badges,
		Enrollments:       enrollments,
	}
}

func fromRepoCampaignSummary(r *scitizenrepo.CampaignSummary) CampaignSummary {
	return CampaignSummary{
		ID:              r.ID,
		Status:          r.Status,
		WindowStart:     r.WindowStart,
		WindowEnd:       r.WindowEnd,
		EnrollmentCount: r.EnrollmentCount,
		RequiredSensors: r.RequiredSensors,
		CreatedAt:       r.CreatedAt,
	}
}

func fromRepoCampaignDetail(r *scitizenrepo.CampaignDetail) *CampaignDetail {
	params := make([]Parameter, len(r.Parameters))
	for i, p := range r.Parameters {
		params[i] = Parameter{Name: p.Name, Unit: p.Unit, MinRange: p.MinRange, MaxRange: p.MaxRange, Precision: p.Precision}
	}
	regions := make([]Region, len(r.Regions))
	for i, rg := range r.Regions {
		regions[i] = Region{GeoJSON: rg.GeoJSON}
	}
	elig := make([]EligibilityCriteria, len(r.Eligibility))
	for i, e := range r.Eligibility {
		elig[i] = EligibilityCriteria{
			DeviceClass: e.DeviceClass, Tier: e.Tier,
			RequiredSensors: e.RequiredSensors, FirmwareMin: e.FirmwareMin,
		}
	}
	return &CampaignDetail{
		CampaignID:      r.CampaignID,
		Status:          r.Status,
		WindowStart:     r.WindowStart,
		WindowEnd:       r.WindowEnd,
		Parameters:      params,
		Regions:         regions,
		Eligibility:     elig,
		EnrollmentCount: r.EnrollmentCount,
		ProgressPercent: r.ProgressPercent,
	}
}

func fromRepoDeviceDetail(r *scitizenrepo.DeviceDetail) *DeviceDetail {
	enrollments := make([]Enrollment, len(r.Enrollments))
	for i, e := range r.Enrollments {
		enrollments[i] = Enrollment{
			ID: e.ID, DeviceID: e.DeviceID, CampaignID: e.CampaignID,
			Status: e.Status, EnrolledAt: e.EnrolledAt,
		}
	}
	connHistory := make([]ConnectionEvent, len(r.ConnectionHistory))
	for i, c := range r.ConnectionHistory {
		connHistory[i] = ConnectionEvent{EventType: c.EventType, Timestamp: c.Timestamp, Reason: c.Reason}
	}
	return &DeviceDetail{
		ID: r.ID, OwnerID: r.OwnerID, Status: r.Status, Class: r.Class,
		FirmwareVersion: r.FirmwareVersion, Tier: r.Tier, Sensors: r.Sensors,
		CertSerial: r.CertSerial, CreatedAt: r.CreatedAt,
		Enrollments: enrollments, ConnectionHistory: connHistory,
	}
}
