package user

// RegisterUserInput is what callers send to RegisterUserFlow.
type RegisterUserInput struct {
	IdpID    string
	UserType string
}

// LoginInput is what callers send to LoginFlow.
type LoginInput struct {
	Email    string
	Password string
}

// LogoutInput is what callers send to LogoutFlow.
type LogoutInput struct {
	SessionID    string
	SessionToken string
}

// RegisterResearcherInput is what callers send to RegisterResearcherFlow.
type RegisterResearcherInput struct {
	Email      string
	Password   string
	GivenName  string
	FamilyName string
}

// VerifyEmailInput is what callers send to VerifyEmailFlow.
type VerifyEmailInput struct {
	UserID           string
	VerificationCode string
}
