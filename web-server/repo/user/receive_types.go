package user

// CreateUserInput is what the CreateUser op sends to the repository.
type CreateUserInput struct {
	IdpID    string
	UserType string
}
