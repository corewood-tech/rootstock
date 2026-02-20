package user

// RegisterUserInput is what callers send to RegisterUserFlow.
type RegisterUserInput struct {
	IdpID    string
	UserType string
}
