package user

import "time"

// User is the user record returned by user flows.
type User struct {
	ID        string
	IdpID     string
	UserType  string
	Status    string
	CreatedAt time.Time
}

// LoginResult is the result of a successful login flow.
type LoginResult struct {
	SessionID    string
	SessionToken string
	User         User
}

// RegisterResearcherResult is the result of researcher registration.
type RegisterResearcherResult struct {
	UserID                 string
	EmailVerificationSent  bool
}
