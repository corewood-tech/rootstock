package user

import (
	"context"

	sessionrepo "rootstock/web-server/repo/session"
	userrepo "rootstock/web-server/repo/user"
)

// Ops holds user operations. Each method is one op.
type Ops struct {
	repo        userrepo.Repository
	sessionRepo sessionrepo.Repository
}

// NewOps creates user ops backed by the given repositories.
func NewOps(repo userrepo.Repository, sessionRepo sessionrepo.Repository) *Ops {
	return &Ops{repo: repo, sessionRepo: sessionRepo}
}

// CreateUser creates a new app user.
func (o *Ops) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	result, err := o.repo.Create(ctx, userrepo.CreateUserInput{
		IdpID:    input.IdpID,
		UserType: input.UserType,
	})
	if err != nil {
		return nil, err
	}
	return fromRepoUser(result), nil
}

// GetUser retrieves a user by ID.
func (o *Ops) GetUser(ctx context.Context, id string) (*User, error) {
	result, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromRepoUser(result), nil
}

// GetUserByIdpID retrieves a user by their identity provider ID.
func (o *Ops) GetUserByIdpID(ctx context.Context, idpID string) (*User, error) {
	result, err := o.repo.GetByIdpID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return fromRepoUser(result), nil
}

// Login creates a session for the given credentials.
func (o *Ops) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	sess, err := o.sessionRepo.CreateSession(ctx, sessionrepo.CreateSessionInput{
		LoginName: input.Email,
		Password:  input.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		SessionID:    sess.SessionID,
		SessionToken: sess.SessionToken,
		UserID:       sess.UserID,
	}, nil
}

// Logout deletes the given session.
func (o *Ops) Logout(ctx context.Context, input LogoutInput) error {
	return o.sessionRepo.DeleteSession(ctx, input.SessionID, input.SessionToken)
}

// ValidateSession verifies a session and returns the associated user ID.
func (o *Ops) ValidateSession(ctx context.Context, input ValidateSessionInput) (*ValidatedSession, error) {
	sess, err := o.sessionRepo.GetSession(ctx, sessionrepo.GetSessionInput{
		SessionID:    input.SessionID,
		SessionToken: input.SessionToken,
	})
	if err != nil {
		return nil, err
	}
	return &ValidatedSession{
		UserID: sess.UserID,
	}, nil
}

// UpdateUserType changes a user's type (researcher, scitizen, both).
func (o *Ops) UpdateUserType(ctx context.Context, input UpdateUserTypeInput) (*User, error) {
	result, err := o.repo.UpdateUserType(ctx, userrepo.UpdateUserTypeInput{
		ID:       input.ID,
		UserType: input.UserType,
	})
	if err != nil {
		return nil, err
	}
	return fromRepoUser(result), nil
}

func fromRepoUser(r *userrepo.User) *User {
	return &User{
		ID:        r.ID,
		IdpID:     r.IdpID,
		UserType:  r.UserType,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
	}
}
