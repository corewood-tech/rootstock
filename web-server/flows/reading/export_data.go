package reading

import (
	"context"

	"rootstock/web-server/ops/pure"
	readingops "rootstock/web-server/ops/reading"
)

// ExportDataFlow orchestrates querying accepted readings and pseudonymizing device IDs.
type ExportDataFlow struct {
	readingOps *readingops.Ops
}

// NewExportDataFlow creates the flow with its required ops.
func NewExportDataFlow(readingOps *readingops.Ops) *ExportDataFlow {
	return &ExportDataFlow{readingOps: readingOps}
}

// Run queries accepted readings for a campaign and pseudonymizes device IDs.
func (f *ExportDataFlow) Run(ctx context.Context, input ExportDataInput) (*ExportDataResult, error) {
	// 1. Query accepted readings
	readings, err := f.readingOps.QueryReadings(ctx, readingops.QueryReadingsInput{
		CampaignID: input.CampaignID,
		Status:     "accepted",
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, err
	}

	// 2. Map to pseudonymizable readings
	pseudoInput := make([]pure.PseudonymizableReading, len(readings))
	for i, r := range readings {
		pseudoInput[i] = pure.PseudonymizableReading{
			DeviceID:        r.DeviceID,
			CampaignID:      r.CampaignID,
			Value:           r.Value,
			Timestamp:       r.Timestamp,
			Geolocation:     r.Geolocation,
			FirmwareVersion: r.FirmwareVersion,
			IngestedAt:      r.IngestedAt,
			Status:          r.Status,
		}
	}

	// 3. Pseudonymize
	pseudonymized := pure.PseudonymizeExport(pure.PseudonymizeInput{
		Readings: pseudoInput,
		Secret:   input.Secret,
	})

	// 4. Map to export result
	exported := make([]ExportedReading, len(pseudonymized))
	for i, p := range pseudonymized {
		exported[i] = ExportedReading{
			PseudoDeviceID:  p.PseudoDeviceID,
			CampaignID:      p.CampaignID,
			Value:           p.Value,
			Timestamp:       p.Timestamp,
			Geolocation:     p.Geolocation,
			FirmwareVersion: p.FirmwareVersion,
			IngestedAt:      p.IngestedAt,
			Status:          p.Status,
		}
	}

	return &ExportDataResult{Readings: exported}, nil
}
