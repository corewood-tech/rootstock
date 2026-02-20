package notification

import (
	"context"

	notificationrepo "rootstock/web-server/repo/notification"
)

// Ops holds notification operations. Each method is one op.
type Ops struct {
	repo notificationrepo.Repository
}

// NewOps creates notification ops backed by the given repository.
func NewOps(repo notificationrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// NotifyScitizens sends notifications to affected scitizens.
// Op #30: FR-033
func (o *Ops) NotifyScitizens(ctx context.Context, input NotifyInput) error {
	notifications := make([]notificationrepo.Notification, len(input.Recipients))
	for i, r := range input.Recipients {
		notifications[i] = notificationrepo.Notification{
			RecipientID: r.ID,
			Subject:     r.Subject,
			Body:        r.Body,
		}
	}
	return o.repo.Notify(ctx, notifications)
}
