package org

import (
	"context"

	orgops "rootstock/web-server/ops/org"
)

// NestOrgFlow creates a child organization under a parent.
type NestOrgFlow struct {
	orgOps *orgops.Ops
}

// NewNestOrgFlow creates the flow with its required ops.
func NewNestOrgFlow(orgOps *orgops.Ops) *NestOrgFlow {
	return &NestOrgFlow{orgOps: orgOps}
}

// Run creates a nested organization.
func (f *NestOrgFlow) Run(ctx context.Context, input NestOrgInput) (*Org, error) {
	result, err := f.orgOps.NestOrg(ctx, orgops.NestOrgInput{
		Name:        input.Name,
		ParentOrgID: input.ParentOrgID,
	})
	if err != nil {
		return nil, err
	}
	return &Org{ID: result.ID, Name: result.Name}, nil
}
