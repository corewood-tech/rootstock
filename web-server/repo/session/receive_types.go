package session

// CreateSessionInput is what callers send to create a session.
type CreateSessionInput struct {
	LoginName string
	Password  string
}

// GetSessionInput is what callers send to verify a session.
type GetSessionInput struct {
	SessionID    string
	SessionToken string
}
