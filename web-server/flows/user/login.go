package user

import (
	"context"
	"fmt"

	userops "rootstock/web-server/ops/user"
)

// LoginFlow orchestrates user login.
// Creates a session, resolves the app user, and returns both.
type LoginFlow struct {
	userOps *userops.Ops
}

// NewLoginFlow creates the flow with its required ops.
func NewLoginFlow(userOps *userops.Ops) *LoginFlow {
	return &LoginFlow{userOps: userOps}
}

// Run authenticates the user and returns session tokens + app user record.
func (f *LoginFlow) Run(ctx context.Context, input LoginInput) (*LoginResult, error) {
	loginResult, err := f.userOps.Login(ctx, userops.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	// Look up app user by Zitadel user ID.
	appUser, err := f.userOps.GetUserByIdpID(ctx, loginResult.UserID)
	if err != nil {
		return nil, err
	}

	// No app record means the user didn't complete registration.
	if appUser == nil {
		return nil, fmt.Errorf("no app record for user %s: registration incomplete", loginResult.UserID)
	}

	return &LoginResult{
		SessionID:    loginResult.SessionID,
		SessionToken: loginResult.SessionToken,
		User:         *fromOpsUser(appUser),
	}, nil
}
