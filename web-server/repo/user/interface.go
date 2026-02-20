package user

import "context"

// Repository defines the interface for app_users data operations.
type Repository interface {
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByIdpID(ctx context.Context, idpID string) (*User, error)
	Shutdown()
}
