package session

import "context"

// Repository defines the interface for session management operations.
type Repository interface {
	CreateSession(ctx context.Context, input CreateSessionInput) (*Session, error)
	GetSession(ctx context.Context, input GetSessionInput) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string, sessionToken string) error
	Shutdown()
}
