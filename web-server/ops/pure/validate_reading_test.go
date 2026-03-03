package pure

import (
	"testing"
	"time"
)

func TestValidateReadingValid(t *testing.T) {
	now := time.Now().UTC()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)
	min := 0.0
	max := 100.0

	result := ValidateReading(
		ReadingInput{Values: map[string]float64{"temp": 23.5}, Timestamp: now},
		ValidationRules{
			WindowStart: &start,
			WindowEnd:   &end,
			Parameters:  []ParameterRule{{Name: "temp", MinRange: &min, MaxRange: &max}},
		},
	)
	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Reason)
	}
}

func TestValidateReadingBeforeWindow(t *testing.T) {
	now := time.Now().UTC()
	start := now.Add(1 * time.Hour) // window starts in the future

	result := ValidateReading(
		ReadingInput{Values: map[string]float64{"temp": 23.5}, Timestamp: now},
		ValidationRules{WindowStart: &start},
	)
	if result.Valid {
		t.Error("expected invalid for timestamp before window")
	}
}

func TestValidateReadingAfterWindow(t *testing.T) {
	now := time.Now().UTC()
	end := now.Add(-1 * time.Hour) // window already ended

	result := ValidateReading(
		ReadingInput{Values: map[string]float64{"temp": 23.5}, Timestamp: now},
		ValidationRules{WindowEnd: &end},
	)
	if result.Valid {
		t.Error("expected invalid for timestamp after window")
	}
}

func TestValidateReadingBelowRange(t *testing.T) {
	min := 10.0
	result := ValidateReading(
		ReadingInput{Values: map[string]float64{"temp": 5.0}, Timestamp: time.Now().UTC()},
		ValidationRules{
			Parameters: []ParameterRule{{Name: "temp", MinRange: &min}},
		},
	)
	if result.Valid {
		t.Error("expected invalid for value below min range")
	}
}

func TestValidateReadingAboveRange(t *testing.T) {
	max := 50.0
	result := ValidateReading(
		ReadingInput{Values: map[string]float64{"temp": 999.0}, Timestamp: time.Now().UTC()},
		ValidationRules{
			Parameters: []ParameterRule{{Name: "temp", MaxRange: &max}},
		},
	)
	if result.Valid {
		t.Error("expected invalid for value above max range")
	}
}

func TestValidateReadingNoRules(t *testing.T) {
	result := ValidateReading(
		ReadingInput{Values: map[string]float64{"temp": 42.0}, Timestamp: time.Now().UTC()},
		ValidationRules{},
	)
	if !result.Valid {
		t.Errorf("expected valid with no rules, got: %s", result.Reason)
	}
}

func TestValidateReadingMultiValue(t *testing.T) {
	minPM := 0.0
	maxPM := 500.0
	minTemp := -40.0
	maxTemp := 85.0

	result := ValidateReading(
		ReadingInput{
			Values:    map[string]float64{"PM2.5": 23.5, "temp": 22.1},
			Timestamp: time.Now().UTC(),
		},
		ValidationRules{
			Parameters: []ParameterRule{
				{Name: "PM2.5", MinRange: &minPM, MaxRange: &maxPM},
				{Name: "temp", MinRange: &minTemp, MaxRange: &maxTemp},
			},
		},
	)
	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Reason)
	}
	if len(result.PerParameter) != 2 {
		t.Errorf("expected 2 per-parameter results, got %d", len(result.PerParameter))
	}
}

func TestValidateReadingMultiValuePartialFail(t *testing.T) {
	maxTemp := 50.0
	minPM := 0.0
	maxPM := 500.0

	result := ValidateReading(
		ReadingInput{
			Values:    map[string]float64{"PM2.5": 23.5, "temp": 999.0},
			Timestamp: time.Now().UTC(),
		},
		ValidationRules{
			Parameters: []ParameterRule{
				{Name: "PM2.5", MinRange: &minPM, MaxRange: &maxPM},
				{Name: "temp", MaxRange: &maxTemp},
			},
		},
	)
	if result.Valid {
		t.Error("expected invalid when one parameter fails")
	}

	// Verify per-parameter results
	var pmValid, tempValid bool
	for _, pv := range result.PerParameter {
		if pv.Name == "PM2.5" {
			pmValid = pv.Valid
		}
		if pv.Name == "temp" {
			tempValid = pv.Valid
		}
	}
	if !pmValid {
		t.Error("expected PM2.5 to be valid")
	}
	if tempValid {
		t.Error("expected temp to be invalid")
	}
}
