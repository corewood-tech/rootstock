package user

// CreateUserInput is what callers send to CreateUser.
type CreateUserInput struct {
	IdpID    string
	UserType string
}
