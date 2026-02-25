package user

import (
	"context"
	"fmt"

	"rootstock/web-server/auth"
	userops "rootstock/web-server/ops/user"
)

// UpdateUserTypeFlow orchestrates changing a user's type.
// Validates the new type and delegates to user ops.
type UpdateUserTypeFlow struct {
	userOps *userops.Ops
}

// NewUpdateUserTypeFlow creates the flow with its required ops.
func NewUpdateUserTypeFlow(userOps *userops.Ops) *UpdateUserTypeFlow {
	return &UpdateUserTypeFlow{userOps: userOps}
}

// Run validates and applies the user type change.
func (f *UpdateUserTypeFlow) Run(ctx context.Context, input UpdateUserTypeInput) (*User, error) {
	if !auth.IsValidUserType(input.UserType) {
		return nil, fmt.Errorf("invalid user_type: %s", input.UserType)
	}

	result, err := f.userOps.UpdateUserType(ctx, userops.UpdateUserTypeInput{
		ID:       input.AppUserID,
		UserType: input.UserType,
	})
	if err != nil {
		return nil, err
	}
	return fromOpsUser(result), nil
}
