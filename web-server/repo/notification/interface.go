package notification

import "context"

// Notification is a single notification to be delivered.
type Notification struct {
	RecipientID string
	Subject     string
	Body        string
}

// Repository defines the interface for notification delivery.
type Repository interface {
	Notify(ctx context.Context, notifications []Notification) error
	Shutdown()
}
