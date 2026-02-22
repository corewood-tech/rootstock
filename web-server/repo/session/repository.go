package session

import (
	"context"
	"fmt"
	"time"

	"github.com/zitadel/zitadel-go/v3/pkg/client"
	sessionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/session/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"google.golang.org/protobuf/types/known/durationpb"

	"rootstock/web-server/config"
)

const sessionLifetime = 7 * 24 * time.Hour

type zitadelRepository struct {
	client *client.Client
}

// NewRepository creates a session repository backed by Zitadel.
func NewRepository(ctx context.Context, cfg config.ZitadelConfig) (Repository, error) {
	z := zitadel.New(
		cfg.Host,
		zitadel.WithInsecure(fmt.Sprintf("%d", cfg.Port)),
	)

	c, err := client.New(ctx, z, client.WithAuth(client.PAT(cfg.ServiceUserPAT)))
	if err != nil {
		return nil, fmt.Errorf("create zitadel client: %w", err)
	}

	return &zitadelRepository{client: c}, nil
}

func (r *zitadelRepository) CreateSession(ctx context.Context, input CreateSessionInput) (*Session, error) {
	resp, err := r.client.SessionServiceV2().CreateSession(ctx, &sessionv2.CreateSessionRequest{
		Checks: &sessionv2.Checks{
			User: &sessionv2.CheckUser{
				Search: &sessionv2.CheckUser_LoginName{LoginName: input.LoginName},
			},
			Password: &sessionv2.CheckPassword{
				Password: input.Password,
			},
		},
		Lifetime: durationpb.New(sessionLifetime),
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	// Retrieve the session to get the user ID from factors.
	getResp, err := r.client.SessionServiceV2().GetSession(ctx, &sessionv2.GetSessionRequest{
		SessionId:    resp.GetSessionId(),
		SessionToken: &resp.SessionToken,
	})
	if err != nil {
		return nil, fmt.Errorf("get created session: %w", err)
	}

	userID := ""
	if s := getResp.GetSession(); s != nil && s.GetFactors() != nil && s.GetFactors().GetUser() != nil {
		userID = s.GetFactors().GetUser().GetId()
	}

	return &Session{
		SessionID:    resp.GetSessionId(),
		SessionToken: resp.GetSessionToken(),
		UserID:       userID,
	}, nil
}

func (r *zitadelRepository) GetSession(ctx context.Context, input GetSessionInput) (*Session, error) {
	resp, err := r.client.SessionServiceV2().GetSession(ctx, &sessionv2.GetSessionRequest{
		SessionId:    input.SessionID,
		SessionToken: &input.SessionToken,
	})
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	s := resp.GetSession()
	if s == nil {
		return nil, fmt.Errorf("session not found")
	}

	userID := ""
	if s.GetFactors() != nil && s.GetFactors().GetUser() != nil {
		userID = s.GetFactors().GetUser().GetId()
	}

	return &Session{
		SessionID:    s.GetId(),
		SessionToken: input.SessionToken,
		UserID:       userID,
	}, nil
}

func (r *zitadelRepository) DeleteSession(ctx context.Context, sessionID string, sessionToken string) error {
	_, err := r.client.SessionServiceV2().DeleteSession(ctx, &sessionv2.DeleteSessionRequest{
		SessionId:    sessionID,
		SessionToken: &sessionToken,
	})
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *zitadelRepository) Shutdown() {
	r.client.Close()
}
