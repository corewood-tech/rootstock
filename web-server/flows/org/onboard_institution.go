package org

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// OnboardInstitutionFlow orchestrates institution onboarding:
// CreateOrg → NestOrg (optional) → DefineRole → AssignRole → InviteUser.
type OnboardInstitutionFlow struct {
	orgOps *orgops.Ops
}

// NewOnboardInstitutionFlow creates the flow with its required ops.
func NewOnboardInstitutionFlow(orgOps *orgops.Ops) *OnboardInstitutionFlow {
	return &OnboardInstitutionFlow{orgOps: orgOps}
}

// Run onboards an institution by creating the org, optionally nesting it,
// defining a default role, assigning it, and inviting the admin user.
func (f *OnboardInstitutionFlow) Run(ctx context.Context, input OnboardInstitutionInput) (*OnboardInstitutionResult, error) {
	// 1. Create the organization
	org, err := f.orgOps.CreateOrg(ctx, orgops.CreateOrgInput{Name: input.OrgName})
	if err != nil {
		return nil, err
	}

	// 2. Nest under parent if specified
	if input.ParentOrgID != "" {
		_, err := f.orgOps.NestOrg(ctx, orgops.NestOrgInput{
			Name:        input.OrgName,
			ParentOrgID: input.ParentOrgID,
		})
		if err != nil {
			return nil, err
		}
	}

	// 3. Define the admin role on the project
	_, err = f.orgOps.DefineRole(ctx, orgops.DefineRoleInput{
		ProjectID:   input.ProjectID,
		RoleKey:     "org_admin",
		DisplayName: "Organization Admin",
	})
	if err != nil {
		return nil, err
	}

	// 4. Invite the admin user
	invite, err := f.orgOps.InviteUser(ctx, orgops.InviteUserInput{
		OrgID:      org.ID,
		Email:      input.AdminEmail,
		GivenName:  input.AdminGivenName,
		FamilyName: input.AdminFamilyName,
	})
	if err != nil {
		return nil, err
	}

	// 5. Assign the admin role to the invited user
	_, err = f.orgOps.AssignRole(ctx, orgops.AssignRoleInput{
		UserID:    invite.UserID,
		ProjectID: input.ProjectID,
		RoleKeys:  []string{"org_admin"},
	})
	if err != nil {
		return nil, err
	}

	return &OnboardInstitutionResult{
		OrgID:  org.ID,
		UserID: invite.UserID,
	}, nil
}
