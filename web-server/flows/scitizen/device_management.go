package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// DeviceManagementFlow handles device listing and detail for scitizens.
// Graph node: 0x27 — implements FR-041 (0x1d), FR-093 (0x1e)
type DeviceManagementFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewDeviceManagementFlow creates the flow with its required ops.
func NewDeviceManagementFlow(scitizenOps *scitizenops.Ops) *DeviceManagementFlow {
	return &DeviceManagementFlow{scitizenOps: scitizenOps}
}

// RunList returns all devices owned by the scitizen.
func (f *DeviceManagementFlow) RunList(ctx context.Context, ownerID string) ([]DeviceSummary, error) {
	results, err := f.scitizenOps.GetDevices(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	out := make([]DeviceSummary, len(results))
	for i, r := range results {
		out[i] = DeviceSummary{
			ID: r.ID, Status: r.Status, Class: r.Class,
			FirmwareVersion: r.FirmwareVersion, Tier: r.Tier, Sensors: r.Sensors,
			ActiveEnrollments: r.ActiveEnrollments, LastSeen: r.LastSeen,
		}
	}
	return out, nil
}

// RunDetail returns full device info with enrollments and connection history.
func (f *DeviceManagementFlow) RunDetail(ctx context.Context, deviceID string) (*DeviceDetail, error) {
	result, err := f.scitizenOps.GetDeviceDetail(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	enrollments := make([]Enrollment, len(result.Enrollments))
	for i, e := range result.Enrollments {
		enrollments[i] = Enrollment{
			ID: e.ID, DeviceID: e.DeviceID, CampaignID: e.CampaignID,
			Status: e.Status, EnrolledAt: e.EnrolledAt,
		}
	}
	connHistory := make([]ConnectionEvent, len(result.ConnectionHistory))
	for i, c := range result.ConnectionHistory {
		connHistory[i] = ConnectionEvent{EventType: c.EventType, Timestamp: c.Timestamp, Reason: c.Reason}
	}

	return &DeviceDetail{
		ID: result.ID, OwnerID: result.OwnerID, Status: result.Status,
		Class: result.Class, FirmwareVersion: result.FirmwareVersion,
		Tier: result.Tier, Sensors: result.Sensors, CertSerial: result.CertSerial,
		CreatedAt: result.CreatedAt, Enrollments: enrollments, ConnectionHistory: connHistory,
	}, nil
}
