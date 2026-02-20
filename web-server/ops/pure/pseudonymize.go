package pure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// PseudonymizeInput is the input to PseudonymizeExport.
type PseudonymizeInput struct {
	Readings []PseudonymizableReading
	Secret   string
}

// PseudonymizableReading is a reading with a real device ID.
type PseudonymizableReading struct {
	DeviceID        string
	CampaignID      string
	Value           float64
	Timestamp       time.Time
	Geolocation     *string
	FirmwareVersion string
	IngestedAt      time.Time
	Status          string
}

// PseudonymizedReading is a reading with the device ID replaced by an HMAC pseudonym.
type PseudonymizedReading struct {
	PseudoDeviceID  string
	CampaignID      string
	Value           float64
	Timestamp       time.Time
	Geolocation     *string
	FirmwareVersion string
	IngestedAt      time.Time
	Status          string
}

// PseudonymizeExport replaces device IDs with HMAC-SHA256 pseudonyms.
// Pure function: no I/O, deterministic for a given secret.
func PseudonymizeExport(input PseudonymizeInput) []PseudonymizedReading {
	out := make([]PseudonymizedReading, len(input.Readings))
	for i, r := range input.Readings {
		mac := hmac.New(sha256.New, []byte(input.Secret))
		mac.Write([]byte(r.DeviceID))
		out[i] = PseudonymizedReading{
			PseudoDeviceID:  hex.EncodeToString(mac.Sum(nil)),
			CampaignID:      r.CampaignID,
			Value:           r.Value,
			Timestamp:       r.Timestamp,
			Geolocation:     r.Geolocation,
			FirmwareVersion: r.FirmwareVersion,
			IngestedAt:      r.IngestedAt,
			Status:          r.Status,
		}
	}
	return out
}
