package pure

import "fmt"

// DeviceCapabilities describes what a device can do.
type DeviceCapabilities struct {
	Class           string
	Tier            int
	Sensors         []string
	FirmwareVersion string
}

// EligibilityCriteria describes what a campaign requires.
type EligibilityCriteria struct {
	DeviceClass     string
	Tier            int
	RequiredSensors []string
	FirmwareMin     string
}

// EligibilityResult is the outcome of the eligibility check.
type EligibilityResult struct {
	Eligible bool
	Reason   string
}

// MatchEligibility is a pure function: (device capabilities, campaign criteria) -> eligible/reason.
func MatchEligibility(caps DeviceCapabilities, criteria EligibilityCriteria) EligibilityResult {
	if caps.Class != criteria.DeviceClass {
		return EligibilityResult{Eligible: false, Reason: fmt.Sprintf("device class %q does not match required %q", caps.Class, criteria.DeviceClass)}
	}

	if caps.Tier < criteria.Tier {
		return EligibilityResult{Eligible: false, Reason: fmt.Sprintf("device tier %d below required %d", caps.Tier, criteria.Tier)}
	}

	if criteria.FirmwareMin != "" && caps.FirmwareVersion < criteria.FirmwareMin {
		return EligibilityResult{Eligible: false, Reason: fmt.Sprintf("firmware %s below minimum %s", caps.FirmwareVersion, criteria.FirmwareMin)}
	}

	sensorSet := make(map[string]bool, len(caps.Sensors))
	for _, s := range caps.Sensors {
		sensorSet[s] = true
	}
	for _, req := range criteria.RequiredSensors {
		if !sensorSet[req] {
			return EligibilityResult{Eligible: false, Reason: fmt.Sprintf("missing required sensor %q", req)}
		}
	}

	return EligibilityResult{Eligible: true, Reason: "eligible"}
}
