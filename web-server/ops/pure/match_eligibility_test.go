package pure

import "testing"

func TestMatchEligibilityPass(t *testing.T) {
	result := MatchEligibility(
		DeviceCapabilities{Class: "weather-station", Tier: 2, Sensors: []string{"temp", "humidity"}, FirmwareVersion: "2.0.0"},
		EligibilityCriteria{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"temp"}, FirmwareMin: "1.0.0"},
	)
	if !result.Eligible {
		t.Errorf("expected eligible, got: %s", result.Reason)
	}
}

func TestMatchEligibilityWrongClass(t *testing.T) {
	result := MatchEligibility(
		DeviceCapabilities{Class: "air-quality", Tier: 1, Sensors: []string{"pm25"}, FirmwareVersion: "1.0.0"},
		EligibilityCriteria{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"temp"}, FirmwareMin: "1.0.0"},
	)
	if result.Eligible {
		t.Error("expected ineligible for wrong class")
	}
}

func TestMatchEligibilityLowTier(t *testing.T) {
	result := MatchEligibility(
		DeviceCapabilities{Class: "weather-station", Tier: 1, Sensors: []string{"temp"}, FirmwareVersion: "1.0.0"},
		EligibilityCriteria{DeviceClass: "weather-station", Tier: 2, RequiredSensors: []string{"temp"}, FirmwareMin: "1.0.0"},
	)
	if result.Eligible {
		t.Error("expected ineligible for low tier")
	}
}

func TestMatchEligibilityOldFirmware(t *testing.T) {
	result := MatchEligibility(
		DeviceCapabilities{Class: "weather-station", Tier: 1, Sensors: []string{"temp"}, FirmwareVersion: "0.9.0"},
		EligibilityCriteria{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"temp"}, FirmwareMin: "1.0.0"},
	)
	if result.Eligible {
		t.Error("expected ineligible for old firmware")
	}
}

func TestMatchEligibilityMissingSensor(t *testing.T) {
	result := MatchEligibility(
		DeviceCapabilities{Class: "weather-station", Tier: 1, Sensors: []string{"temp"}, FirmwareVersion: "1.0.0"},
		EligibilityCriteria{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"temp", "humidity"}, FirmwareMin: "1.0.0"},
	)
	if result.Eligible {
		t.Error("expected ineligible for missing sensor")
	}
}

func TestMatchEligibilityNoFirmwareMin(t *testing.T) {
	result := MatchEligibility(
		DeviceCapabilities{Class: "weather-station", Tier: 1, Sensors: []string{"temp"}, FirmwareVersion: "0.1.0"},
		EligibilityCriteria{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"temp"}, FirmwareMin: ""},
	)
	if !result.Eligible {
		t.Errorf("expected eligible with no firmware min, got: %s", result.Reason)
	}
}
