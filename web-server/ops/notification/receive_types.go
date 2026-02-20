package notification

// NotifyInput is what callers send to NotifyScitizens.
type NotifyInput struct {
	Recipients []Recipient
}

// Recipient is a single notification target.
type Recipient struct {
	ID      string
	Subject string
	Body    string
}
