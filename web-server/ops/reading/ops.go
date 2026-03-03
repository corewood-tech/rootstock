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

// QuarantineReadingValue flags an individual reading value as quarantined.
func (o *Ops) QuarantineReadingValue(ctx context.Context, readingValueID string, reason string) error {
	return o.repo.QuarantineValue(ctx, readingValueID, reason)
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
	qm := &QualityMetrics{
		CampaignID:      result.CampaignID,
		AcceptedCount:   result.AcceptedCount,
		QuarantineCount: result.QuarantineCount,
	}
	for _, pq := range result.PerParameter {
		qm.PerParameter = append(qm.PerParameter, ParameterQuality{
			ParameterName:    pq.ParameterName,
			AcceptedCount:    pq.AcceptedCount,
			QuarantinedCount: pq.QuarantinedCount,
		})
	}
	return qm, nil
}

// GetCampaignDeviceBreakdown returns per-device stats with pseudonymized IDs.
func (o *Ops) GetCampaignDeviceBreakdown(ctx context.Context, campaignID string, hmacSecret string) ([]DeviceBreakdown, error) {
	results, err := o.repo.GetCampaignDeviceBreakdown(ctx, campaignID, hmacSecret)
	if err != nil {
		return nil, err
	}
	out := make([]DeviceBreakdown, len(results))
	for i, r := range results {
		out[i] = DeviceBreakdown{
			PseudoDeviceID: r.PseudoDeviceID,
			DeviceClass:    r.DeviceClass,
			AcceptanceRate: r.AcceptanceRate,
			ReadingCount:   r.ReadingCount,
			LastSeen:       r.LastSeen,
		}
	}
	return out, nil
}

// GetCampaignTemporalCoverage returns hourly reading counts.
func (o *Ops) GetCampaignTemporalCoverage(ctx context.Context, campaignID string) ([]TemporalBucket, error) {
	results, err := o.repo.GetCampaignTemporalCoverage(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	out := make([]TemporalBucket, len(results))
	for i, r := range results {
		out[i] = TemporalBucket{Bucket: r.Bucket, Count: r.Count}
	}
	return out, nil
}

// GetEnrollmentFunnel returns enrollment stage counts.
func (o *Ops) GetEnrollmentFunnel(ctx context.Context, campaignID string) (*EnrollmentFunnel, error) {
	result, err := o.repo.GetEnrollmentFunnel(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return &EnrollmentFunnel{
		Enrolled:     result.Enrolled,
		Active:       result.Active,
		Contributing: result.Contributing,
	}, nil
}

// GetScitizenReadingStats returns aggregated reading stats for a scitizen across all their devices.
func (o *Ops) GetScitizenReadingStats(ctx context.Context, scitizenID string) (*ScitizenReadingStats, error) {
	result, err := o.repo.GetScitizenReadingStats(ctx, scitizenID)
	if err != nil {
		return nil, err
	}
	return &ScitizenReadingStats{
		Volume:      result.Volume,
		QualityRate: result.QualityRate,
		Consistency: result.Consistency,
		Diversity:   result.Diversity,
	}, nil
}

func toRepoPersistInput(in PersistReadingInput) readingrepo.PersistReadingInput {
	values := make([]readingrepo.ReadingValueInput, len(in.Values))
	for i, v := range in.Values {
		values[i] = readingrepo.ReadingValueInput{
			ParameterName: v.ParameterName,
			Value:         v.Value,
		}
	}
	return readingrepo.PersistReadingInput{
		DeviceID:        in.DeviceID,
		CampaignID:      in.CampaignID,
		Values:          values,
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
		Offset:     in.Offset,
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
	rd := &Reading{
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
	for _, rv := range r.Values {
		rd.Values = append(rd.Values, ReadingValue{
			ID:               rv.ID,
			ReadingID:        rv.ReadingID,
			ParameterName:    rv.ParameterName,
			Value:            rv.Value,
			Status:           rv.Status,
			QuarantineReason: rv.QuarantineReason,
		})
	}
	return rd
}
