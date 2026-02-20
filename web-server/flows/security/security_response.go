package security

import (
	"context"
	"fmt"

	deviceops "rootstock/web-server/ops/device"
	notificationops "rootstock/web-server/ops/notification"
	readingops "rootstock/web-server/ops/reading"
)

// SecurityResponseFlow suspends vulnerable devices by class/firmware,
// quarantines readings from the vulnerability window, and notifies affected scitizens.
type SecurityResponseFlow struct {
	deviceOps       *deviceops.Ops
	readingOps      *readingops.Ops
	notificationOps *notificationops.Ops
}

// NewSecurityResponseFlow creates the flow with its required ops.
func NewSecurityResponseFlow(deviceOps *deviceops.Ops, readingOps *readingops.Ops, notificationOps *notificationops.Ops) *SecurityResponseFlow {
	return &SecurityResponseFlow{
		deviceOps:       deviceOps,
		readingOps:      readingOps,
		notificationOps: notificationOps,
	}
}

// Run executes the security response: suspend, quarantine, notify.
func (f *SecurityResponseFlow) Run(ctx context.Context, input SecurityResponseInput) (*SecurityResponseResult, error) {
	// 1. Query affected devices by class and firmware range.
	devices, err := f.deviceOps.QueryDevicesByClass(ctx, deviceops.QueryByClassInput{
		Class:          input.Class,
		FirmwareMinGte: input.FirmwareMin,
		FirmwareMaxLte: input.FirmwareMax,
	})
	if err != nil {
		return nil, fmt.Errorf("query devices by class: %w", err)
	}

	// 2. Suspend each device.
	for _, d := range devices {
		if err := f.deviceOps.UpdateDeviceStatus(ctx, d.ID, "suspended"); err != nil {
			return nil, fmt.Errorf("suspend device %s: %w", d.ID, err)
		}
	}

	// 3. Collect device IDs and quarantine readings from the vulnerability window.
	deviceIDs := make([]string, len(devices))
	for i, d := range devices {
		deviceIDs[i] = d.ID
	}

	var quarantined int64
	if len(deviceIDs) > 0 {
		quarantined, err = f.readingOps.QuarantineByWindow(ctx, readingops.QuarantineByWindowInput{
			DeviceIDs: deviceIDs,
			Since:     input.WindowStart,
			Until:     input.WindowEnd,
			Reason:    input.Reason,
		})
		if err != nil {
			return nil, fmt.Errorf("quarantine readings: %w", err)
		}
	}

	// 4. Deduplicate owner IDs and notify affected scitizens.
	seen := make(map[string]bool)
	var recipients []notificationops.Recipient
	for _, d := range devices {
		if seen[d.OwnerID] {
			continue
		}
		seen[d.OwnerID] = true
		recipients = append(recipients, notificationops.Recipient{
			ID:      d.OwnerID,
			Subject: "Device suspended due to security vulnerability",
			Body:    fmt.Sprintf("Your device (class %s) has been suspended: %s. Affected firmware range: %s–%s.", input.Class, input.Reason, input.FirmwareMin, input.FirmwareMax),
		})
	}

	if len(recipients) > 0 {
		if err := f.notificationOps.NotifyScitizens(ctx, notificationops.NotifyInput{
			Recipients: recipients,
		}); err != nil {
			return nil, fmt.Errorf("notify scitizens: %w", err)
		}
	}

	return &SecurityResponseResult{
		SuspendedCount:      len(devices),
		QuarantinedReadings: quarantined,
		NotifiedScitizens:   len(recipients),
	}, nil
}
