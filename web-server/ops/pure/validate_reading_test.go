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
		ReadingInput{Value: 23.5, Timestamp: now},
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
		ReadingInput{Value: 23.5, Timestamp: now},
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
		ReadingInput{Value: 23.5, Timestamp: now},
		ValidationRules{WindowEnd: &end},
	)
	if result.Valid {
		t.Error("expected invalid for timestamp after window")
	}
}

func TestValidateReadingBelowRange(t *testing.T) {
	min := 10.0
	result := ValidateReading(
		ReadingInput{Value: 5.0, Timestamp: time.Now().UTC()},
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
		ReadingInput{Value: 999.0, Timestamp: time.Now().UTC()},
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
		ReadingInput{Value: 42.0, Timestamp: time.Now().UTC()},
		ValidationRules{},
	)
	if !result.Valid {
		t.Errorf("expected valid with no rules, got: %s", result.Reason)
	}
}
