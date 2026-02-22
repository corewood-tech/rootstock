package user

import (
	"context"

	userops "rootstock/web-server/ops/user"
)

// LogoutFlow orchestrates user logout by deleting the session.
type LogoutFlow struct {
	userOps *userops.Ops
}

// NewLogoutFlow creates the flow with its required ops.
func NewLogoutFlow(userOps *userops.Ops) *LogoutFlow {
	return &LogoutFlow{userOps: userOps}
}

// Run deletes the session.
func (f *LogoutFlow) Run(ctx context.Context, input LogoutInput) error {
	return f.userOps.Logout(ctx, userops.LogoutInput{
		SessionID:    input.SessionID,
		SessionToken: input.SessionToken,
	})
}
