package notification

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
)

type notifyReq struct {
	ctx           context.Context
	notifications []Notification
	resp          chan error
}

type shutdownReq struct {
	resp chan struct{}
}

type smtpRepo struct {
	host       string
	port       int
	from       string
	notifyCh   chan notifyReq
	shutdownCh chan shutdownReq
}

// NewRepository creates an SMTP-based notification repository.
func NewRepository(host string, port int, from string) Repository {
	r := &smtpRepo{
		host:       host,
		port:       port,
		from:       from,
		notifyCh:   make(chan notifyReq),
		shutdownCh: make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *smtpRepo) manage() {
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

func (r *smtpRepo) Notify(ctx context.Context, notifications []Notification) error {
	resp := make(chan error, 1)
	r.notifyCh <- notifyReq{ctx: ctx, notifications: notifications, resp: resp}
	return <-resp
}

func (r *smtpRepo) Shutdown() {
	resp := make(chan struct{})
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

func (r *smtpRepo) doNotify(ctx context.Context, notifications []Notification) error {
	addr := fmt.Sprintf("%s:%d", r.host, r.port)

	for _, n := range notifications {
		msg := strings.Join([]string{
			"From: " + r.from,
			"To: " + n.RecipientID,
			"Subject: " + n.Subject,
			"MIME-Version: 1.0",
			"Content-Type: text/plain; charset=utf-8",
			"",
			n.Body,
		}, "\r\n")

		if err := smtp.SendMail(addr, nil, r.from, []string{n.RecipientID}, []byte(msg)); err != nil {
			slog.ErrorContext(ctx, "smtp send failed",
				"recipient", n.RecipientID,
				"subject", n.Subject,
				"error", err,
			)
			return fmt.Errorf("send notification to %s: %w", n.RecipientID, err)
		}

		slog.InfoContext(ctx, "notification sent",
			"recipient", n.RecipientID,
			"subject", n.Subject,
		)
	}
	return nil
}
