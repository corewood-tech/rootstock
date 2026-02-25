package user

import (
	"context"
	"fmt"

	notificationops "rootstock/web-server/ops/notification"
	orgops "rootstock/web-server/ops/org"
	userops "rootstock/web-server/ops/user"
)

// RegisterResearcherFlow orchestrates researcher self-registration.
// Creates a Zitadel user (with password) and an app record.
// Uses ReturnCode to get the verification code, then sends the
// verification email via notification ops (app controls delivery).
type RegisterResearcherFlow struct {
	orgOps            *orgops.Ops
	userOps           *userops.Ops
	notificationOps   *notificationops.Ops
	verifyURLBase     string
}

// NewRegisterResearcherFlow creates the flow with its required ops.
func NewRegisterResearcherFlow(orgOps *orgops.Ops, userOps *userops.Ops, notificationOps *notificationops.Ops, verifyURLBase string) *RegisterResearcherFlow {
	return &RegisterResearcherFlow{
		orgOps:          orgOps,
		userOps:         userOps,
		notificationOps: notificationOps,
		verifyURLBase:   verifyURLBase,
	}
}

// Run creates the Zitadel user, the app record, and sends the verification email.
func (f *RegisterResearcherFlow) Run(ctx context.Context, input RegisterResearcherInput) (*RegisterResearcherResult, error) {
	// 1. Create user in Zitadel with password. ReturnCode gives us the verification code.
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
	_, err = f.userOps.CreateUser(ctx, userops.CreateUserInput{
		IdpID:    idpUser.UserID,
		UserType: input.UserType,
	})
	if err != nil {
		return nil, err
	}

	// 3. Send verification email via notification ops.
	verifyURL := fmt.Sprintf("%s/app/en/verify-email?userId=%s&code=%s",
		f.verifyURLBase, idpUser.UserID, idpUser.EmailCode)

	body := fmt.Sprintf(
		"Welcome to Rootstock!\n\nPlease verify your email by clicking the link below:\n\n%s\n\nIf you did not create this account, you can ignore this email.",
		verifyURL,
	)

	err = f.notificationOps.NotifyScitizens(ctx, notificationops.NotifyInput{
		Recipients: []notificationops.Recipient{
			{
				ID:      input.Email,
				Subject: "Verify your Rootstock account",
				Body:    body,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send verification email: %w", err)
	}

	return &RegisterResearcherResult{
		UserID:                idpUser.UserID,
		EmailVerificationSent: true,
	}, nil
}
