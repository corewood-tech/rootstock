package reading

import (
	"context"

	readingrepo "rootstock/web-server/repo/reading"
)

// Ops holds reading operations. Each method is one op.
type Ops struct {
	repo readingrepo.Repository
}

// NewOps creates reading ops backed by the given repository.
func NewOps(repo readingrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// PersistReading writes a reading with full provenance.
// Op #20: FR-023
func (o *Ops) PersistReading(ctx context.Context, input PersistReadingInput) (*Reading, error) {
	result, err := o.repo.Persist(ctx, toRepoPersistInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoReading(result), nil
}

// QuarantineReading flags a reading as quarantined.
// Op #21: FR-025
func (o *Ops) QuarantineReading(ctx context.Context, id string, reason string) error {
	return o.repo.Quarantine(ctx, id, reason)
}

// QueryReadings reads campaign data with provenance.
// Op #22: FR-026
func (o *Ops) QueryReadings(ctx context.Context, input QueryReadingsInput) ([]Reading, error) {
	results, err := o.repo.Query(ctx, toRepoQueryInput(input))
	if err != nil {
		return nil, err
	}
	out := make([]Reading, len(results))
	for i, r := range results {
		out[i] = *fromRepoReading(&r)
	}
	return out, nil
}

// QuarantineByWindow batch-flags readings from affected devices during a vulnerability window.
// Op #24: FR-032
func (o *Ops) QuarantineByWindow(ctx context.Context, input QuarantineByWindowInput) (int64, error) {
	return o.repo.QuarantineByWindow(ctx, toRepoQuarantineByWindowInput(input))
}

// GetCampaignQuality returns aggregated metrics for a campaign.
// Op #25: FR-010
func (o *Ops) GetCampaignQuality(ctx context.Context, campaignID string) (*QualityMetrics, error) {
	result, err := o.repo.GetCampaignQuality(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return &QualityMetrics{
		CampaignID:      result.CampaignID,
		AcceptedCount:   result.AcceptedCount,
		QuarantineCount: result.QuarantineCount,
	}, nil
}

func toRepoPersistInput(in PersistReadingInput) readingrepo.PersistReadingInput {
	return readingrepo.PersistReadingInput{
		DeviceID:        in.DeviceID,
		CampaignID:      in.CampaignID,
		Value:           in.Value,
		Timestamp:       in.Timestamp,
		Geolocation:     in.Geolocation,
		FirmwareVersion: in.FirmwareVersion,
		CertSerial:      in.CertSerial,
	}
}

func toRepoQueryInput(in QueryReadingsInput) readingrepo.QueryReadingsInput {
	return readingrepo.QueryReadingsInput{
		CampaignID: in.CampaignID,
		DeviceID:   in.DeviceID,
		Status:     in.Status,
		Since:      in.Since,
		Until:      in.Until,
		Limit:      in.Limit,
	}
}

func toRepoQuarantineByWindowInput(in QuarantineByWindowInput) readingrepo.QuarantineByWindowInput {
	return readingrepo.QuarantineByWindowInput{
		DeviceIDs: in.DeviceIDs,
		Since:     in.Since,
		Until:     in.Until,
		Reason:    in.Reason,
	}
}

func fromRepoReading(r *readingrepo.Reading) *Reading {
	return &Reading{
		ID:               r.ID,
		DeviceID:         r.DeviceID,
		CampaignID:       r.CampaignID,
		Value:            r.Value,
		Timestamp:        r.Timestamp,
		Geolocation:      r.Geolocation,
		FirmwareVersion:  r.FirmwareVersion,
		CertSerial:       r.CertSerial,
		IngestedAt:       r.IngestedAt,
		Status:           r.Status,
		QuarantineReason: r.QuarantineReason,
	}
}
