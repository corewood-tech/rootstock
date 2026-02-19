package org

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// CreateOrgFlow creates an organization.
type CreateOrgFlow struct {
	orgOps *orgops.Ops
}

// NewCreateOrgFlow creates the flow with its required ops.
func NewCreateOrgFlow(orgOps *orgops.Ops) *CreateOrgFlow {
	return &CreateOrgFlow{orgOps: orgOps}
}

// Run creates an organization.
func (f *CreateOrgFlow) Run(ctx context.Context, input CreateOrgInput) (*Org, error) {
	result, err := f.orgOps.CreateOrg(ctx, orgops.CreateOrgInput{Name: input.Name})
	if err != nil {
		return nil, err
	}
	return &Org{ID: result.ID, Name: result.Name}, nil
}
