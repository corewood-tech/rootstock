package user

import (
	"context"

	userrepo "rootstock/web-server/repo/user"
)

// Ops holds user operations. Each method is one op.
type Ops struct {
	repo userrepo.Repository
}

// NewOps creates user ops backed by the given repository.
func NewOps(repo userrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// CreateUser creates a new app user.
func (o *Ops) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	result, err := o.repo.Create(ctx, userrepo.CreateUserInput{
		IdpID:    input.IdpID,
		UserType: input.UserType,
	})
	if err != nil {
		return nil, err
	}
	return fromRepoUser(result), nil
}

// GetUser retrieves a user by ID.
func (o *Ops) GetUser(ctx context.Context, id string) (*User, error) {
	result, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromRepoUser(result), nil
}

// GetUserByIdpID retrieves a user by their identity provider ID.
func (o *Ops) GetUserByIdpID(ctx context.Context, idpID string) (*User, error) {
	result, err := o.repo.GetByIdpID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return fromRepoUser(result), nil
}

func fromRepoUser(r *userrepo.User) *User {
	return &User{
		ID:        r.ID,
		IdpID:     r.IdpID,
		UserType:  r.UserType,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
	}
}
