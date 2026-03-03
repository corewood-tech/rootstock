package pure

import (
	"fmt"
	"time"
)

// ReadingInput is the reading data to validate. Supports multi-value readings.
type ReadingInput struct {
	Values      map[string]float64 // parameter name -> value
	Timestamp   time.Time
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

// ParameterValidation is the validation result for a single parameter.
type ParameterValidation struct {
	Name   string
	Valid  bool
	Reason string
}

// ValidationResult is the outcome of reading validation.
type ValidationResult struct {
	Valid        bool
	Reason       string
	PerParameter []ParameterValidation
}

// ValidateReading is a pure function: (reading, rules) -> valid/invalid + reason.
// Each parameter is validated independently. Overall Valid = all parameters pass + timestamp valid.
func ValidateReading(input ReadingInput, rules ValidationRules) ValidationResult {
	// Check timestamp within campaign window
	if rules.WindowStart != nil && input.Timestamp.Before(*rules.WindowStart) {
		return ValidationResult{Valid: false, Reason: fmt.Sprintf("timestamp %s before campaign window start %s", input.Timestamp.Format(time.RFC3339), rules.WindowStart.Format(time.RFC3339))}
	}
	if rules.WindowEnd != nil && input.Timestamp.After(*rules.WindowEnd) {
		return ValidationResult{Valid: false, Reason: fmt.Sprintf("timestamp %s after campaign window end %s", input.Timestamp.Format(time.RFC3339), rules.WindowEnd.Format(time.RFC3339))}
	}

	// Build a rules lookup by parameter name
	rulesByName := make(map[string]ParameterRule, len(rules.Parameters))
	for _, p := range rules.Parameters {
		rulesByName[p.Name] = p
	}

	allValid := true
	var perParam []ParameterValidation

	for name, value := range input.Values {
		pv := ParameterValidation{Name: name, Valid: true, Reason: "valid"}
		if rule, ok := rulesByName[name]; ok {
			if rule.MinRange != nil && value < *rule.MinRange {
				pv.Valid = false
				pv.Reason = fmt.Sprintf("value %f below min range %f for %s", value, *rule.MinRange, name)
				allValid = false
			}
			if pv.Valid && rule.MaxRange != nil && value > *rule.MaxRange {
				pv.Valid = false
				pv.Reason = fmt.Sprintf("value %f above max range %f for %s", value, *rule.MaxRange, name)
				allValid = false
			}
		}
		perParam = append(perParam, pv)
	}

	reason := "valid"
	if !allValid {
		for _, pv := range perParam {
			if !pv.Valid {
				reason = pv.Reason
				break
			}
		}
	}

	return ValidationResult{
		Valid:        allValid,
		Reason:       reason,
		PerParameter: perParam,
	}
}
