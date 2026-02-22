package user

import (
	"context"
	"fmt"
	"log/slog"

	notificationops "rootstock/web-server/ops/notification"
	orgops "rootstock/web-server/ops/org"
	userops "rootstock/web-server/ops/user"
)

// RegisterResearcherFlow orchestrates researcher self-registration.
// Creates a Zitadel user (with password), an app record, and a session.
type RegisterResearcherFlow struct {
	orgOps          *orgops.Ops
	userOps         *userops.Ops
	notificationOps *notificationops.Ops
}

// NewRegisterResearcherFlow creates the flow with its required ops.
func NewRegisterResearcherFlow(orgOps *orgops.Ops, userOps *userops.Ops, notificationOps *notificationops.Ops) *RegisterResearcherFlow {
	return &RegisterResearcherFlow{orgOps: orgOps, userOps: userOps, notificationOps: notificationOps}
}

// Run creates the Zitadel user, the app record, and a session.
func (f *RegisterResearcherFlow) Run(ctx context.Context, input RegisterResearcherInput) (*RegisterResearcherResult, error) {
	// 1. Create user in Zitadel with password.
	idpUser, err := f.orgOps.CreateIdpUser(ctx, orgops.CreateIdpUserInput{
		Email:      input.Email,
		Password:   input.Password,
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
	})
	if err != nil {
		return nil, err
	}

	// 2. Create app record.
	appUser, err := f.userOps.CreateUser(ctx, userops.CreateUserInput{
		IdpID:    idpUser.UserID,
		UserType: "researcher",
	})
	if err != nil {
		return nil, err
	}

	// 3. Create session (log them in immediately).
	loginResult, err := f.userOps.Login(ctx, userops.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	// 4. Send welcome email (best-effort).
	if err := f.notificationOps.NotifyScitizens(ctx, notificationops.NotifyInput{
		Recipients: []notificationops.Recipient{{
			ID:      input.Email,
			Subject: "Welcome to Rootstock",
			Body:    fmt.Sprintf("Hello %s,\n\nYour researcher account has been created.\n\nWelcome to Rootstock by Corewood.", input.GivenName),
		}},
	}); err != nil {
		slog.WarnContext(ctx, "failed to send welcome email", "email", input.Email, "error", err)
	}

	return &RegisterResearcherResult{
		SessionID:    loginResult.SessionID,
		SessionToken: loginResult.SessionToken,
		User:         *fromOpsUser(appUser),
	}, nil
}
