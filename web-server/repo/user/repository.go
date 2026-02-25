package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type response[T any] struct {
	val T
	err error
}

type createReq struct {
	ctx   context.Context
	input CreateUserInput
	resp  chan response[*User]
}

type getByIDReq struct {
	ctx  context.Context
	id   string
	resp chan response[*User]
}

type getByIdpIDReq struct {
	ctx   context.Context
	idpID string
	resp  chan response[*User]
}

type updateUserTypeReq struct {
	ctx   context.Context
	input UpdateUserTypeInput
	resp  chan response[*User]
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool             *pgxpool.Pool
	createCh         chan createReq
	getByIDCh        chan getByIDReq
	getByIdpIDCh     chan getByIdpIDReq
	updateUserTypeCh chan updateUserTypeReq
	shutdownCh       chan shutdownReq
}

// NewRepository creates a user repository backed by Postgres.
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:             pool,
		createCh:         make(chan createReq),
		getByIDCh:        make(chan getByIDReq),
		getByIdpIDCh:     make(chan getByIdpIDReq),
		updateUserTypeCh: make(chan updateUserTypeReq),
		shutdownCh:       make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *pgRepo) manage() {
	for {
		select {
		case req := <-r.createCh:
			val, err := r.doCreate(req.ctx, req.input)
			req.resp <- response[*User]{val: val, err: err}
		case req := <-r.getByIDCh:
			val, err := r.doGetByID(req.ctx, req.id)
			req.resp <- response[*User]{val: val, err: err}
		case req := <-r.getByIdpIDCh:
			val, err := r.doGetByIdpID(req.ctx, req.idpID)
			req.resp <- response[*User]{val: val, err: err}
		case req := <-r.updateUserTypeCh:
			val, err := r.doUpdateUserType(req.ctx, req.input)
			req.resp <- response[*User]{val: val, err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	resp := make(chan response[*User], 1)
	r.createCh <- createReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetByID(ctx context.Context, id string) (*User, error) {
	resp := make(chan response[*User], 1)
	r.getByIDCh <- getByIDReq{ctx: ctx, id: id, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetByIdpID(ctx context.Context, idpID string) (*User, error) {
	resp := make(chan response[*User], 1)
	r.getByIdpIDCh <- getByIdpIDReq{ctx: ctx, idpID: idpID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) UpdateUserType(ctx context.Context, input UpdateUserTypeInput) (*User, error) {
	resp := make(chan response[*User], 1)
	r.updateUserTypeCh <- updateUserTypeReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doCreate(ctx context.Context, input CreateUserInput) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO app_users (id, idp_id, user_type)
		 VALUES ($1, $2, $3)
		 RETURNING id, idp_id, user_type, status, created_at`,
		ulid.Make().String(), input.IdpID, input.UserType,
	).Scan(&u.ID, &u.IdpID, &u.UserType, &u.Status, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &u, nil
}

func (r *pgRepo) doGetByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`SELECT id, idp_id, user_type, status, created_at FROM app_users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.IdpID, &u.UserType, &u.Status, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user %s not found", id)
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, nil
}

func (r *pgRepo) doUpdateUserType(ctx context.Context, input UpdateUserTypeInput) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`UPDATE app_users SET user_type = $2
		 WHERE id = $1
		 RETURNING id, idp_id, user_type, status, created_at`,
		input.ID, input.UserType,
	).Scan(&u.ID, &u.IdpID, &u.UserType, &u.Status, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user %s not found", input.ID)
		}
		return nil, fmt.Errorf("update user type: %w", err)
	}
	return &u, nil
}

func (r *pgRepo) doGetByIdpID(ctx context.Context, idpID string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`SELECT id, idp_id, user_type, status, created_at FROM app_users WHERE idp_id = $1`,
		idpID,
	).Scan(&u.ID, &u.IdpID, &u.UserType, &u.Status, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by idp_id: %w", err)
	}
	return &u, nil
}
