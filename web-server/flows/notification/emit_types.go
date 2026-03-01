package notification

import "time"

// ListResult is the result of a notification list query.
type ListResult struct {
	Notifications []Notification
	UnreadCount   int
	Total         int
}

// Notification is an in-app notification record.
type Notification struct {
	ID           string
	UserID       string
	Type         string
	Message      string
	Read         bool
	ResourceLink *string
	CreatedAt    time.Time
}

// MarkReadResult is the result of marking notifications as read.
type MarkReadResult struct {
	MarkedCount int
}

// Preference is a notification preference entry.
type Preference struct {
	Type  string
	InApp bool
	Email bool
}
