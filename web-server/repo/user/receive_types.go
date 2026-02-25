package user

// CreateUserInput is what the CreateUser op sends to the repository.
type CreateUserInput struct {
	IdpID    string
	UserType string
}

// UpdateUserTypeInput is what the UpdateUserType op sends to the repository.
type UpdateUserTypeInput struct {
	ID       string
	UserType string
}
