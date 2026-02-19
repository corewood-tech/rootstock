package pure

import (
	"fmt"
	"time"
)

// ReadingInput is the reading data to validate.
type ReadingInput struct {
	Value     float64
	Timestamp time.Time
	Geolocation *GeoPoint
}

// GeoPoint is a lat/lon pair.
type GeoPoint struct {
	Longitude float64
	Latitude  float64
}

// ValidationRules are the campaign rules to validate against.
type ValidationRules struct {
	Parameters  []ParameterRule
	WindowStart *time.Time
	WindowEnd   *time.Time
}

// ParameterRule defines valid ranges for a measurement parameter.
type ParameterRule struct {
	Name     string
	MinRange *float64
	MaxRange *float64
}

// ValidationResult is the outcome of reading validation.
type ValidationResult struct {
	Valid  bool
	Reason string
}

// ValidateReading is a pure function: (reading, rules) -> valid/invalid + reason.
func ValidateReading(input ReadingInput, rules ValidationRules) ValidationResult {
	// Check timestamp within campaign window
	if rules.WindowStart != nil && input.Timestamp.Before(*rules.WindowStart) {
		return ValidationResult{Valid: false, Reason: fmt.Sprintf("timestamp %s before campaign window start %s", input.Timestamp.Format(time.RFC3339), rules.WindowStart.Format(time.RFC3339))}
	}
	if rules.WindowEnd != nil && input.Timestamp.After(*rules.WindowEnd) {
		return ValidationResult{Valid: false, Reason: fmt.Sprintf("timestamp %s after campaign window end %s", input.Timestamp.Format(time.RFC3339), rules.WindowEnd.Format(time.RFC3339))}
	}

	// Check value against parameter ranges
	for _, p := range rules.Parameters {
		if p.MinRange != nil && input.Value < *p.MinRange {
			return ValidationResult{Valid: false, Reason: fmt.Sprintf("value %f below min range %f for %s", input.Value, *p.MinRange, p.Name)}
		}
		if p.MaxRange != nil && input.Value > *p.MaxRange {
			return ValidationResult{Valid: false, Reason: fmt.Sprintf("value %f above max range %f for %s", input.Value, *p.MaxRange, p.Name)}
		}
	}

	return ValidationResult{Valid: true, Reason: "valid"}
}
