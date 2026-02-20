package device

import (
	"context"
	"encoding/json"

	campaignops "rootstock/web-server/ops/campaign"
	deviceops "rootstock/web-server/ops/device"
	mqttops "rootstock/web-server/ops/mqtt"
	"rootstock/web-server/ops/pure"
)

// EnrollInCampaignFlow orchestrates device-to-campaign enrollment:
// GetDeviceCapabilities → GetCampaignEligibility → MatchEligibility →
// EnrollDeviceInCampaign → PushDeviceConfig.
type EnrollInCampaignFlow struct {
	deviceOps   *deviceops.Ops
	campaignOps *campaignops.Ops
	mqttOps     *mqttops.Ops
}

// NewEnrollInCampaignFlow creates the flow with its required ops.
func NewEnrollInCampaignFlow(deviceOps *deviceops.Ops, campaignOps *campaignops.Ops, mqttOps *mqttops.Ops) *EnrollInCampaignFlow {
	return &EnrollInCampaignFlow{
		deviceOps:   deviceOps,
		campaignOps: campaignOps,
		mqttOps:     mqttOps,
	}
}

// Run checks eligibility, enrolls the device, and pushes config via MQTT.
func (f *EnrollInCampaignFlow) Run(ctx context.Context, input EnrollInCampaignInput) (*EnrollInCampaignResult, error) {
	// 1. Get device capabilities
	caps, err := f.deviceOps.GetDeviceCapabilities(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}

	// 2. Get campaign eligibility criteria
	criteria, err := f.campaignOps.GetCampaignEligibility(ctx, input.CampaignID)
	if err != nil {
		return nil, err
	}

	// 3. Check eligibility against each criterion (any match = eligible)
	var lastResult pure.EligibilityResult
	eligible := false
	for _, c := range criteria {
		lastResult = pure.MatchEligibility(
			pure.DeviceCapabilities{
				Class:           caps.Class,
				Tier:            caps.Tier,
				Sensors:         caps.Sensors,
				FirmwareVersion: caps.FirmwareVersion,
			},
			pure.EligibilityCriteria{
				DeviceClass:     c.DeviceClass,
				Tier:            c.Tier,
				RequiredSensors: c.RequiredSensors,
				FirmwareMin:     c.FirmwareMin,
			},
		)
		if lastResult.Eligible {
			eligible = true
			break
		}
	}

	if !eligible {
		reason := "no eligibility criteria defined"
		if len(criteria) > 0 {
			reason = lastResult.Reason
		}
		return &EnrollInCampaignResult{Enrolled: false, Reason: reason}, nil
	}

	// 4. Enroll device in campaign
	if err := f.deviceOps.EnrollDeviceInCampaign(ctx, input.DeviceID, input.CampaignID); err != nil {
		return nil, err
	}

	// 5. Push config to device via MQTT
	configPayload, err := json.Marshal(DeviceConfigPayload{
		CampaignID: input.CampaignID,
	})
	if err != nil {
		return nil, err
	}

	if err := f.mqttOps.PushDeviceConfig(ctx, mqttops.PushDeviceConfigInput{
		DeviceID: input.DeviceID,
		Payload:  configPayload,
	}); err != nil {
		return nil, err
	}

	return &EnrollInCampaignResult{Enrolled: true}, nil
}
