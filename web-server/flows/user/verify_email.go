package user

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// VerifyEmailFlow verifies a user's email address using the code from the verification link.
type VerifyEmailFlow struct {
	orgOps *orgops.Ops
}

// NewVerifyEmailFlow creates the flow with its required ops.
func NewVerifyEmailFlow(orgOps *orgops.Ops) *VerifyEmailFlow {
	return &VerifyEmailFlow{orgOps: orgOps}
}

// Run verifies the user's email using the verification code from Zitadel.
func (f *VerifyEmailFlow) Run(ctx context.Context, input VerifyEmailInput) error {
	return f.orgOps.VerifyEmail(ctx, orgops.VerifyEmailInput{
		UserID:           input.UserID,
		VerificationCode: input.VerificationCode,
	})
}
