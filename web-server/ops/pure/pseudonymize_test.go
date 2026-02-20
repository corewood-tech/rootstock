package pure

import (
	"testing"
	"time"
)

func TestPseudonymizeExport(t *testing.T) {
	now := time.Now()
	geo := "point"

	input := PseudonymizeInput{
		Readings: []PseudonymizableReading{
			{
				DeviceID:        "dev-1",
				CampaignID:      "camp-1",
				Value:           22.5,
				Timestamp:       now,
				Geolocation:     &geo,
				FirmwareVersion: "1.0.0",
				IngestedAt:      now,
				Status:          "accepted",
			},
			{
				DeviceID:        "dev-2",
				CampaignID:      "camp-1",
				Value:           23.1,
				Timestamp:       now,
				Geolocation:     nil,
				FirmwareVersion: "1.0.0",
				IngestedAt:      now,
				Status:          "accepted",
			},
		},
		Secret: "test-secret",
	}

	results := PseudonymizeExport(input)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Pseudonyms should not be the original device IDs
	if results[0].PseudoDeviceID == "dev-1" {
		t.Error("pseudonym should not equal original device ID")
	}
	if results[1].PseudoDeviceID == "dev-2" {
		t.Error("pseudonym should not equal original device ID")
	}

	// Different devices should produce different pseudonyms
	if results[0].PseudoDeviceID == results[1].PseudoDeviceID {
		t.Error("different devices should have different pseudonyms")
	}

	// Same device + secret should be deterministic
	results2 := PseudonymizeExport(input)
	if results[0].PseudoDeviceID != results2[0].PseudoDeviceID {
		t.Error("pseudonymization should be deterministic")
	}

	// Other fields should be preserved
	if results[0].CampaignID != "camp-1" {
		t.Errorf("campaign_id = %s, want camp-1", results[0].CampaignID)
	}
	if results[0].Value != 22.5 {
		t.Errorf("value = %f, want 22.5", results[0].Value)
	}
	if results[0].Geolocation == nil || *results[0].Geolocation != "point" {
		t.Error("geolocation should be preserved")
	}
	if results[1].Geolocation != nil {
		t.Error("nil geolocation should remain nil")
	}

	// Different secret should produce different pseudonyms
	input2 := input
	input2.Secret = "different-secret"
	results3 := PseudonymizeExport(input2)
	if results[0].PseudoDeviceID == results3[0].PseudoDeviceID {
		t.Error("different secrets should produce different pseudonyms")
	}
}

func TestPseudonymizeExportEmpty(t *testing.T) {
	results := PseudonymizeExport(PseudonymizeInput{
		Readings: nil,
		Secret:   "test-secret",
	})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
