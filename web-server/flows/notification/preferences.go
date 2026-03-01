package notification

import (
	"context"

	enrollmentops "rootstock/web-server/ops/enrollment"
)

// GetPreferencesFlow handles retrieving notification preferences.
type GetPreferencesFlow struct {
	enrollmentOps *enrollmentops.Ops
}

// NewGetPreferencesFlow creates the flow with its required ops.
func NewGetPreferencesFlow(enrollmentOps *enrollmentops.Ops) *GetPreferencesFlow {
	return &GetPreferencesFlow{enrollmentOps: enrollmentOps}
}

// Run returns notification preferences for the user.
func (f *GetPreferencesFlow) Run(ctx context.Context, userID string) ([]Preference, error) {
	results, err := f.enrollmentOps.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]Preference, len(results))
	for i, r := range results {
		out[i] = Preference{Type: r.Type, InApp: r.InApp, Email: r.Email}
	}
	return out, nil
}

// UpdatePreferencesFlow handles updating notification preferences.
type UpdatePreferencesFlow struct {
	enrollmentOps *enrollmentops.Ops
}

// NewUpdatePreferencesFlow creates the flow with its required ops.
func NewUpdatePreferencesFlow(enrollmentOps *enrollmentops.Ops) *UpdatePreferencesFlow {
	return &UpdatePreferencesFlow{enrollmentOps: enrollmentOps}
}

// Run updates notification preferences for the user.
func (f *UpdatePreferencesFlow) Run(ctx context.Context, input UpdatePreferencesInput) error {
	prefs := make([]enrollmentops.NotificationPreference, len(input.Preferences))
	for i, p := range input.Preferences {
		prefs[i] = enrollmentops.NotificationPreference{Type: p.Type, InApp: p.InApp, Email: p.Email}
	}
	return f.enrollmentOps.UpdatePreferences(ctx, input.UserID, prefs)
}
