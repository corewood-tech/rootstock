package org

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// AssignRoleFlow grants roles to a user on a project.
type AssignRoleFlow struct {
	orgOps *orgops.Ops
}

// NewAssignRoleFlow creates the flow with its required ops.
func NewAssignRoleFlow(orgOps *orgops.Ops) *AssignRoleFlow {
	return &AssignRoleFlow{orgOps: orgOps}
}

// Run assigns roles to a user.
func (f *AssignRoleFlow) Run(ctx context.Context, input AssignRoleInput) (*UserGrant, error) {
	result, err := f.orgOps.AssignRole(ctx, orgops.AssignRoleInput{
		UserID:    input.UserID,
		ProjectID: input.ProjectID,
		RoleKeys:  input.RoleKeys,
	})
	if err != nil {
		return nil, err
	}
	return &UserGrant{
		UserGrantID: result.UserGrantID,
		UserID:      result.UserID,
		ProjectID:   result.ProjectID,
		RoleKeys:    result.RoleKeys,
	}, nil
}
