package user

import "time"

// User is the user record returned by user ops.
type User struct {
	ID        string
	IdpID     string
	UserType  string
	Status    string
	CreatedAt time.Time
}

// LoginResult is the result of a successful login.
type LoginResult struct {
	SessionID    string
	SessionToken string
	UserID       string
}

// ValidatedSession is the result of a session validation.
type ValidatedSession struct {
	UserID string
}
