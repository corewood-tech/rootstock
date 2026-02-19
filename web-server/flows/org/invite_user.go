package org

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// InviteUserFlow creates a user and sends an invitation.
type InviteUserFlow struct {
	orgOps *orgops.Ops
}

// NewInviteUserFlow creates the flow with its required ops.
func NewInviteUserFlow(orgOps *orgops.Ops) *InviteUserFlow {
	return &InviteUserFlow{orgOps: orgOps}
}

// Run invites a user to an organization.
func (f *InviteUserFlow) Run(ctx context.Context, input InviteUserInput) (*InviteResult, error) {
	result, err := f.orgOps.InviteUser(ctx, orgops.InviteUserInput{
		OrgID:      input.OrgID,
		Email:      input.Email,
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
	})
	if err != nil {
		return nil, err
	}
	return &InviteResult{
		UserID:    result.UserID,
		EmailCode: result.EmailCode,
	}, nil
}
