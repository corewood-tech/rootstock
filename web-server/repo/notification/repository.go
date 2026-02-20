package notification

import (
	"context"
	"log/slog"
)

type notifyReq struct {
	ctx           context.Context
	notifications []Notification
	resp          chan error
}

type shutdownReq struct {
	resp chan struct{}
}

type logRepo struct {
	notifyCh   chan notifyReq
	shutdownCh chan shutdownReq
}

// NewRepository creates a log-based notification repository.
// Real provider (email/push) deferred.
func NewRepository() Repository {
	r := &logRepo{
		notifyCh:   make(chan notifyReq),
		shutdownCh: make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *logRepo) manage() {
	for {
		select {
		case req := <-r.notifyCh:
			req.resp <- r.doNotify(req.ctx, req.notifications)

		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *logRepo) Notify(ctx context.Context, notifications []Notification) error {
	resp := make(chan error, 1)
	r.notifyCh <- notifyReq{ctx: ctx, notifications: notifications, resp: resp}
	return <-resp
}

func (r *logRepo) Shutdown() {
	resp := make(chan struct{})
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

func (r *logRepo) doNotify(ctx context.Context, notifications []Notification) error {
	for _, n := range notifications {
		slog.InfoContext(ctx, "notification sent",
			"recipient_id", n.RecipientID,
			"subject", n.Subject,
			"body", n.Body,
		)
	}
	return nil
}
