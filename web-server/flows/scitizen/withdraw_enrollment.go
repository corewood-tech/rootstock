package scitizen

import (
	"context"

	enrollmentops "rootstock/web-server/ops/enrollment"
)

// WithdrawEnrollmentFlow withdraws a device from a campaign.
// Graph node: 0x2e — implements FR-063 (0x10)
type WithdrawEnrollmentFlow struct {
	enrollmentOps *enrollmentops.Ops
}

// NewWithdrawEnrollmentFlow creates the flow with its required ops.
func NewWithdrawEnrollmentFlow(enrollmentOps *enrollmentops.Ops) *WithdrawEnrollmentFlow {
	return &WithdrawEnrollmentFlow{enrollmentOps: enrollmentOps}
}

// Run withdraws a device enrollment by enrollment ID.
func (f *WithdrawEnrollmentFlow) Run(ctx context.Context, enrollmentID string) error {
	return f.enrollmentOps.Withdraw(ctx, enrollmentID)
}
