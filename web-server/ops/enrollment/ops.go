package enrollment

import (
	"context"

	enrollmentrepo "rootstock/web-server/repo/enrollment"
)

// Ops holds enrollment operations. Each method is one op.
// Graph node: 0x2b (EnrollmentOps)
type Ops struct {
	repo enrollmentrepo.Repository
}

// NewOps creates enrollment ops backed by the given repository.
func NewOps(repo enrollmentrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// Enroll creates a campaign enrollment with consent record.
func (o *Ops) Enroll(ctx context.Context, input EnrollInput) (*Enrollment, error) {
	result, err := o.repo.Enroll(ctx, enrollmentrepo.EnrollInput{
		DeviceID:       input.DeviceID,
		CampaignID:     input.CampaignID,
		ScitizenID:     input.ScitizenID,
		ConsentVersion: input.ConsentVersion,
		ConsentScope:   input.ConsentScope,
	})
	if err != nil {
		return nil, err
	}
	return fromRepoEnrollment(result), nil
}

// Withdraw withdraws a device from a campaign.
func (o *Ops) Withdraw(ctx context.Context, enrollmentID string) error {
	return o.repo.Withdraw(ctx, enrollmentID)
}

// GetByID returns an enrollment by ID.
func (o *Ops) GetByID(ctx context.Context, id string) (*Enrollment, error) {
	result, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromRepoEnrollment(result), nil
}

// GetByDeviceCampaign returns an enrollment by device+campaign.
func (o *Ops) GetByDeviceCampaign(ctx context.Context, deviceID, campaignID string) (*Enrollment, error) {
	result, err := o.repo.GetByDeviceCampaign(ctx, deviceID, campaignID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return fromRepoEnrollment(result), nil
}

// MarkRead marks notifications as read.
func (o *Ops) MarkRead(ctx context.Context, userID string, ids []string) (int, error) {
	return o.repo.MarkRead(ctx, userID, ids)
}

// CreateNotification creates a notification.
func (o *Ops) CreateNotification(ctx context.Context, input CreateNotificationInput) error {
	return o.repo.CreateNotification(ctx, enrollmentrepo.CreateNotificationInput{
		UserID:       input.UserID,
		Type:         input.Type,
		Message:      input.Message,
		ResourceLink: input.ResourceLink,
	})
}

// GetPreferences returns notification preferences.
func (o *Ops) GetPreferences(ctx context.Context, userID string) ([]NotificationPreference, error) {
	results, err := o.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]NotificationPreference, len(results))
	for i, r := range results {
		out[i] = NotificationPreference{Type: r.Type, InApp: r.InApp, Email: r.Email}
	}
	return out, nil
}

// UpdatePreferences updates notification preferences.
func (o *Ops) UpdatePreferences(ctx context.Context, userID string, prefs []NotificationPreference) error {
	repoPrefs := make([]enrollmentrepo.NotificationPreference, len(prefs))
	for i, p := range prefs {
		repoPrefs[i] = enrollmentrepo.NotificationPreference{Type: p.Type, InApp: p.InApp, Email: p.Email}
	}
	return o.repo.UpdatePreferences(ctx, userID, repoPrefs)
}

func fromRepoEnrollment(r *enrollmentrepo.Enrollment) *Enrollment {
	return &Enrollment{
		ID:          r.ID,
		DeviceID:    r.DeviceID,
		CampaignID:  r.CampaignID,
		ScitizenID:  r.ScitizenID,
		Status:      r.Status,
		EnrolledAt:  r.EnrolledAt,
		WithdrawnAt: r.WithdrawnAt,
	}
}
