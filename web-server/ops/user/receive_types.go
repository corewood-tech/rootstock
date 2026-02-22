package user

// CreateUserInput is what callers send to CreateUser.
type CreateUserInput struct {
	IdpID    string
	UserType string
}

// LoginInput is what callers send to Login.
type LoginInput struct {
	Email    string
	Password string
}

// LogoutInput is what callers send to Logout.
type LogoutInput struct {
	SessionID    string
	SessionToken string
}

// ValidateSessionInput is what callers send to ValidateSession.
type ValidateSessionInput struct {
	SessionID    string
	SessionToken string
}
