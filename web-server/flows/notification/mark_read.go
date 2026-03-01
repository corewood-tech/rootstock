package notification

import (
	"context"

	enrollmentops "rootstock/web-server/ops/enrollment"
)

// MarkReadFlow handles marking notifications as read.
type MarkReadFlow struct {
	enrollmentOps *enrollmentops.Ops
}

// NewMarkReadFlow creates the flow with its required ops.
func NewMarkReadFlow(enrollmentOps *enrollmentops.Ops) *MarkReadFlow {
	return &MarkReadFlow{enrollmentOps: enrollmentOps}
}

// Run marks the specified notifications as read and returns the count.
func (f *MarkReadFlow) Run(ctx context.Context, input MarkReadInput) (*MarkReadResult, error) {
	count, err := f.enrollmentOps.MarkRead(ctx, input.UserID, input.NotificationIDs)
	if err != nil {
		return nil, err
	}
	return &MarkReadResult{MarkedCount: count}, nil
}
