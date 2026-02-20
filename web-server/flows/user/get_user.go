package user

import (
	"context"

	userops "rootstock/web-server/ops/user"
)

// GetUserFlow retrieves a user by their idp_id.
type GetUserFlow struct {
	userOps *userops.Ops
}

// NewGetUserFlow creates the flow with its required ops.
func NewGetUserFlow(userOps *userops.Ops) *GetUserFlow {
	return &GetUserFlow{userOps: userOps}
}

// Run retrieves a user by their idp_id (resolved from JWT in the handler).
func (f *GetUserFlow) Run(ctx context.Context, idpID string) (*User, error) {
	result, err := f.userOps.GetUserByIdpID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return fromOpsUser(result), nil
}
