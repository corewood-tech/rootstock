package notification

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// ListNotificationsFlow handles listing notifications for a user.
type ListNotificationsFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewListNotificationsFlow creates the flow with its required ops.
func NewListNotificationsFlow(scitizenOps *scitizenops.Ops) *ListNotificationsFlow {
	return &ListNotificationsFlow{scitizenOps: scitizenOps}
}

// Run returns notifications for the user.
func (f *ListNotificationsFlow) Run(ctx context.Context, input ListInput) (*ListResult, error) {
	results, unreadCount, total, err := f.scitizenOps.GetNotifications(ctx, scitizenops.GetNotificationsInput{
		UserID:     input.UserID,
		TypeFilter: input.TypeFilter,
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, err
	}

	out := make([]Notification, len(results))
	for i, r := range results {
		out[i] = Notification{
			ID: r.ID, UserID: r.UserID, Type: r.Type, Message: r.Message,
			Read: r.Read, ResourceLink: r.ResourceLink, CreatedAt: r.CreatedAt,
		}
	}

	return &ListResult{
		Notifications: out,
		UnreadCount:   unreadCount,
		Total:         total,
	}, nil
}
