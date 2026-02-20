package user

import (
	"context"

	userops "rootstock/web-server/ops/user"
)

// RegisterUserFlow orchestrates user registration.
// Idempotent — if the user already exists by idp_id, returns the existing record.
type RegisterUserFlow struct {
	userOps *userops.Ops
}

// NewRegisterUserFlow creates the flow with its required ops.
func NewRegisterUserFlow(userOps *userops.Ops) *RegisterUserFlow {
	return &RegisterUserFlow{userOps: userOps}
}

// Run checks if a user exists by idp_id. If yes, returns existing. If no, creates new.
func (f *RegisterUserFlow) Run(ctx context.Context, input RegisterUserInput) (*User, error) {
	existing, err := f.userOps.GetUserByIdpID(ctx, input.IdpID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return fromOpsUser(existing), nil
	}

	created, err := f.userOps.CreateUser(ctx, userops.CreateUserInput{
		IdpID:    input.IdpID,
		UserType: input.UserType,
	})
	if err != nil {
		return nil, err
	}
	return fromOpsUser(created), nil
}

func fromOpsUser(u *userops.User) *User {
	return &User{
		ID:        u.ID,
		IdpID:     u.IdpID,
		UserType:  u.UserType,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
