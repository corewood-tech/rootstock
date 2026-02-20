package security

// SecurityResponseResult is the outcome of a security response action.
type SecurityResponseResult struct {
	SuspendedCount      int
	QuarantinedReadings int64
	NotifiedScitizens   int
}
