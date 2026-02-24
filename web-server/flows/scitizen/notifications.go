package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// NotificationFlow handles notification listing for scitizens.
// Graph node: 0x23 — implements FR-114 (0x5), FR-051 (0xc)
type NotificationFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewNotificationFlow creates the flow with its required ops.
func NewNotificationFlow(scitizenOps *scitizenops.Ops) *NotificationFlow {
	return &NotificationFlow{scitizenOps: scitizenOps}
}

// Run returns notifications for the scitizen.
func (f *NotificationFlow) Run(ctx context.Context, input GetNotificationsInput) (*NotificationsResult, error) {
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

	return &NotificationsResult{
		Notifications: out,
		UnreadCount:   unreadCount,
		Total:         total,
	}, nil
}
