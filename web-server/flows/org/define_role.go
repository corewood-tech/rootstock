package org

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// DefineRoleFlow defines a role on a project.
type DefineRoleFlow struct {
	orgOps *orgops.Ops
}

// NewDefineRoleFlow creates the flow with its required ops.
func NewDefineRoleFlow(orgOps *orgops.Ops) *DefineRoleFlow {
	return &DefineRoleFlow{orgOps: orgOps}
}

// Run defines a role.
func (f *DefineRoleFlow) Run(ctx context.Context, input DefineRoleInput) (*Role, error) {
	result, err := f.orgOps.DefineRole(ctx, orgops.DefineRoleInput{
		ProjectID:   input.ProjectID,
		RoleKey:     input.RoleKey,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	return &Role{
		ProjectID:   result.ProjectID,
		RoleKey:     result.RoleKey,
		DisplayName: result.DisplayName,
	}, nil
}
