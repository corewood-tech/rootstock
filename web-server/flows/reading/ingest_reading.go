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

// Run validates a reading against campaign rules, persists it, and quarantines invalid values.
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
			Values:    input.Values,
			Timestamp: input.Timestamp,
		},
		pure.ValidationRules{
			Parameters:  paramRules,
			WindowStart: rules.WindowStart,
			WindowEnd:   rules.WindowEnd,
		},
	)

	// 3. Persist the reading with all values
	opsInput := toOpsReadingInput(input)
	opsReading, err := f.readingOps.PersistReading(ctx, opsInput)
	if err != nil {
		return nil, err
	}

	// 4. If timestamp invalid, quarantine the whole reading
	if !validationResult.Valid && len(validationResult.PerParameter) == 0 {
		if err := f.readingOps.QuarantineReading(ctx, opsReading.ID, validationResult.Reason); err != nil {
			return nil, err
		}
		opsReading.Status = "quarantined"
		opsReading.QuarantineReason = &validationResult.Reason
	}

	// 5. Quarantine individual values that failed validation
	failedParams := make(map[string]string) // name -> reason
	for _, pv := range validationResult.PerParameter {
		if !pv.Valid {
			failedParams[pv.Name] = pv.Reason
		}
	}
	for i := range opsReading.Values {
		if reason, failed := failedParams[opsReading.Values[i].ParameterName]; failed {
			if err := f.readingOps.QuarantineReadingValue(ctx, opsReading.Values[i].ID, reason); err != nil {
				return nil, err
			}
			opsReading.Values[i].Status = "quarantined"
			opsReading.Values[i].QuarantineReason = &reason
		}
	}

	// 6. If all values are quarantined, quarantine the reading itself
	if len(opsReading.Values) > 0 {
		allQuarantined := true
		for _, v := range opsReading.Values {
			if v.Status != "quarantined" {
				allQuarantined = false
				break
			}
		}
		if allQuarantined {
			reason := "all parameter values quarantined"
			if err := f.readingOps.QuarantineReading(ctx, opsReading.ID, reason); err != nil {
				return nil, err
			}
			opsReading.Status = "quarantined"
			opsReading.QuarantineReason = &reason
		}
	}

	// 7. Anomaly detection for accepted values (best-effort, per parameter)
	if opsReading.Status == "accepted" {
		for paramName, value := range input.Values {
			if _, failed := failedParams[paramName]; failed {
				continue
			}

			// Update rolling baseline
			if _, err := f.graphOps.UpdateBaseline(ctx, graphops.UpdateBaselineInput{
				CampaignRef:   input.CampaignID,
				ParameterName: paramName,
				Value:         value,
			}); err != nil {
				slog.WarnContext(ctx, "failed to update baseline", "campaign_id", input.CampaignID, "parameter", paramName, "error", err)
			}

			// Check for anomaly
			anomaly, err := f.graphOps.CheckAnomaly(ctx, graphops.CheckAnomalyInput{
				CampaignRef:   input.CampaignID,
				ParameterName: paramName,
				Value:         value,
			})
			if err != nil {
				slog.WarnContext(ctx, "failed to check anomaly", "campaign_id", input.CampaignID, "parameter", paramName, "error", err)
			} else if anomaly != nil {
				// Find the reading value and quarantine it
				for i := range opsReading.Values {
					if opsReading.Values[i].ParameterName == paramName {
						if err := f.readingOps.QuarantineReadingValue(ctx, opsReading.Values[i].ID, anomaly.Reason); err != nil {
							slog.WarnContext(ctx, "failed to quarantine anomaly value", "reading_value_id", opsReading.Values[i].ID, "error", err)
						} else {
							opsReading.Values[i].Status = "quarantined"
							opsReading.Values[i].QuarantineReason = &anomaly.Reason
						}
						break
					}
				}
			}
		}
	}

	return fromOpsReading(opsReading), nil
}

func toOpsReadingInput(in IngestReadingInput) readingops.PersistReadingInput {
	values := make([]readingops.ReadingValueInput, 0, len(in.Values))
	for name, value := range in.Values {
		values = append(values, readingops.ReadingValueInput{
			ParameterName: name,
			Value:         value,
		})
	}
	return readingops.PersistReadingInput{
		DeviceID:        in.DeviceID,
		CampaignID:      in.CampaignID,
		Values:          values,
		Timestamp:       in.Timestamp,
		Geolocation:     in.Geolocation,
		FirmwareVersion: in.FirmwareVersion,
		CertSerial:      in.CertSerial,
	}
}

func fromOpsReading(r *readingops.Reading) *Reading {
	rd := &Reading{
		ID:               r.ID,
		DeviceID:         r.DeviceID,
		CampaignID:       r.CampaignID,
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
