package reading

import (
	"context"
	"log/slog"

	campaignops "rootstock/web-server/ops/campaign"
	graphops "rootstock/web-server/ops/graph"
	"rootstock/web-server/ops/pure"
	readingops "rootstock/web-server/ops/reading"
)

// IngestReadingFlow orchestrates reading ingestion: validate then persist.
type IngestReadingFlow struct {
	campaignOps *campaignops.Ops
	readingOps  *readingops.Ops
	graphOps    *graphops.Ops
}

// NewIngestReadingFlow creates the flow with its required ops.
func NewIngestReadingFlow(campaignOps *campaignops.Ops, readingOps *readingops.Ops, graphOps *graphops.Ops) *IngestReadingFlow {
	return &IngestReadingFlow{campaignOps: campaignOps, readingOps: readingOps, graphOps: graphOps}
}

// Run validates a reading against campaign rules, persists it, and quarantines if invalid.
func (f *IngestReadingFlow) Run(ctx context.Context, input IngestReadingInput) (*Reading, error) {
	// 1. Get campaign validation rules
	rules, err := f.campaignOps.GetCampaignRules(ctx, input.CampaignID)
	if err != nil {
		return nil, err
	}

	// 2. Validate the reading (pure op — no I/O)
	var paramRules []pure.ParameterRule
	for _, p := range rules.Parameters {
		paramRules = append(paramRules, pure.ParameterRule{
			Name:     p.Name,
			MinRange: p.MinRange,
			MaxRange: p.MaxRange,
		})
	}
	validationResult := pure.ValidateReading(
		pure.ReadingInput{
			Value:     input.Value,
			Timestamp: input.Timestamp,
		},
		pure.ValidationRules{
			Parameters:  paramRules,
			WindowStart: rules.WindowStart,
			WindowEnd:   rules.WindowEnd,
		},
	)

	// 3. Persist the reading
	opsInput := toOpsReadingInput(input)
	opsReading, err := f.readingOps.PersistReading(ctx, opsInput)
	if err != nil {
		return nil, err
	}

	// 4. If invalid, quarantine it
	if !validationResult.Valid {
		if err := f.readingOps.QuarantineReading(ctx, opsReading.ID, validationResult.Reason); err != nil {
			return nil, err
		}
		opsReading.Status = "quarantined"
		opsReading.QuarantineReason = &validationResult.Reason
	}

	// 5. Anomaly detection for valid readings (best-effort)
	if validationResult.Valid && len(rules.Parameters) > 0 {
		paramName := rules.Parameters[0].Name

		// Update rolling baseline
		if _, err := f.graphOps.UpdateBaseline(ctx, graphops.UpdateBaselineInput{
			CampaignRef:   input.CampaignID,
			ParameterName: paramName,
			Value:         input.Value,
		}); err != nil {
			slog.WarnContext(ctx, "failed to update baseline", "campaign_id", input.CampaignID, "error", err)
		}

		// Check for anomaly
		anomaly, err := f.graphOps.CheckAnomaly(ctx, graphops.CheckAnomalyInput{
			CampaignRef:   input.CampaignID,
			ParameterName: paramName,
			Value:         input.Value,
		})
		if err != nil {
			slog.WarnContext(ctx, "failed to check anomaly", "campaign_id", input.CampaignID, "error", err)
		} else if anomaly != nil {
			if err := f.readingOps.QuarantineReading(ctx, opsReading.ID, anomaly.Reason); err != nil {
				slog.WarnContext(ctx, "failed to quarantine anomaly", "reading_id", opsReading.ID, "error", err)
			} else {
				opsReading.Status = "quarantined"
				opsReading.QuarantineReason = &anomaly.Reason
			}
		}
	}

	return fromOpsReading(opsReading), nil
}

func toOpsReadingInput(in IngestReadingInput) readingops.PersistReadingInput {
	return readingops.PersistReadingInput{
		DeviceID:        in.DeviceID,
		CampaignID:      in.CampaignID,
		Value:           in.Value,
		Timestamp:       in.Timestamp,
		Geolocation:     in.Geolocation,
		FirmwareVersion: in.FirmwareVersion,
		CertSerial:      in.CertSerial,
	}
}

func fromOpsReading(r *readingops.Reading) *Reading {
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
