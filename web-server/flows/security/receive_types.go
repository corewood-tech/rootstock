package security

import "time"

// SecurityResponseInput is what callers send to SecurityResponseFlow.Run.
type SecurityResponseInput struct {
	Class       string
	FirmwareMin string
	FirmwareMax string
	WindowStart time.Time
	WindowEnd   time.Time
	Reason      string
}
