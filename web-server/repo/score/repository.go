package score

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

type upsertReq struct {
	ctx   context.Context
	input UpsertScoreInput
	resp  chan response[*Score]
}

type getScoreReq struct {
	ctx        context.Context
	scitizenID string
	resp       chan response[*Score]
}

type awardBadgeReq struct {
	ctx        context.Context
	scitizenID string
	badgeType  string
	resp       chan response[struct{}]
}

type getBadgesReq struct {
	ctx        context.Context
	scitizenID string
	resp       chan response[[]Badge]
}

type grantSweepReq struct {
	ctx   context.Context
	input GrantSweepstakesInput
	resp  chan response[struct{}]
}

type getSweepReq struct {
	ctx        context.Context
	scitizenID string
	resp       chan response[[]SweepstakesEntry]
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool           *pgxpool.Pool
	upsertCh      chan upsertReq
	getScoreCh    chan getScoreReq
	awardBadgeCh  chan awardBadgeReq
	getBadgesCh   chan getBadgesReq
	grantSweepCh  chan grantSweepReq
	getSweepCh    chan getSweepReq
	shutdownCh    chan shutdownReq
}

// NewRepository creates a score repository backed by Postgres.
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:          pool,
		upsertCh:     make(chan upsertReq),
		getScoreCh:   make(chan getScoreReq),
		awardBadgeCh: make(chan awardBadgeReq),
		getBadgesCh:  make(chan getBadgesReq),
		grantSweepCh: make(chan grantSweepReq),
		getSweepCh:   make(chan getSweepReq),
		shutdownCh:   make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *pgRepo) manage() {
	for {
		select {
		case req := <-r.upsertCh:
			val, err := r.doUpsert(req.ctx, req.input)
			req.resp <- response[*Score]{val: val, err: err}
		case req := <-r.getScoreCh:
			val, err := r.doGetScore(req.ctx, req.scitizenID)
			req.resp <- response[*Score]{val: val, err: err}
		case req := <-r.awardBadgeCh:
			err := r.doAwardBadge(req.ctx, req.scitizenID, req.badgeType)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.getBadgesCh:
			val, err := r.doGetBadges(req.ctx, req.scitizenID)
			req.resp <- response[[]Badge]{val: val, err: err}
		case req := <-r.grantSweepCh:
			err := r.doGrantSweepstakes(req.ctx, req.input)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.getSweepCh:
			val, err := r.doGetSweepstakes(req.ctx, req.scitizenID)
			req.resp <- response[[]SweepstakesEntry]{val: val, err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) UpsertScore(ctx context.Context, input UpsertScoreInput) (*Score, error) {
	resp := make(chan response[*Score], 1)
	r.upsertCh <- upsertReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetScore(ctx context.Context, scitizenID string) (*Score, error) {
	resp := make(chan response[*Score], 1)
	r.getScoreCh <- getScoreReq{ctx: ctx, scitizenID: scitizenID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) AwardBadge(ctx context.Context, scitizenID string, badgeType string) error {
	resp := make(chan response[struct{}], 1)
	r.awardBadgeCh <- awardBadgeReq{ctx: ctx, scitizenID: scitizenID, badgeType: badgeType, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) GetBadges(ctx context.Context, scitizenID string) ([]Badge, error) {
	resp := make(chan response[[]Badge], 1)
	r.getBadgesCh <- getBadgesReq{ctx: ctx, scitizenID: scitizenID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GrantSweepstakes(ctx context.Context, input GrantSweepstakesInput) error {
	resp := make(chan response[struct{}], 1)
	r.grantSweepCh <- grantSweepReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) GetSweepstakesEntries(ctx context.Context, scitizenID string) ([]SweepstakesEntry, error) {
	resp := make(chan response[[]SweepstakesEntry], 1)
	r.getSweepCh <- getSweepReq{ctx: ctx, scitizenID: scitizenID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doUpsert(ctx context.Context, input UpsertScoreInput) (*Score, error) {
	var s Score
	err := r.pool.QueryRow(ctx,
		`INSERT INTO scores (scitizen_id, volume, quality_rate, consistency, diversity, total)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (scitizen_id) DO UPDATE SET
			volume = EXCLUDED.volume,
			quality_rate = EXCLUDED.quality_rate,
			consistency = EXCLUDED.consistency,
			diversity = EXCLUDED.diversity,
			total = EXCLUDED.total,
			updated_at = now()
		 RETURNING scitizen_id, volume, quality_rate, consistency, diversity, total, updated_at`,
		input.ScitizenID, input.Volume, input.QualityRate, input.Consistency, input.Diversity, input.Total,
	).Scan(&s.ScitizenID, &s.Volume, &s.QualityRate, &s.Consistency, &s.Diversity, &s.Total, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert score: %w", err)
	}
	return &s, nil
}

func (r *pgRepo) doGetScore(ctx context.Context, scitizenID string) (*Score, error) {
	var s Score
	err := r.pool.QueryRow(ctx,
		`SELECT scitizen_id, volume, quality_rate, consistency, diversity, total, updated_at
		 FROM scores WHERE scitizen_id = $1`,
		scitizenID,
	).Scan(&s.ScitizenID, &s.Volume, &s.QualityRate, &s.Consistency, &s.Diversity, &s.Total, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("score for %s not found", scitizenID)
		}
		return nil, fmt.Errorf("get score: %w", err)
	}
	return &s, nil
}

func (r *pgRepo) doAwardBadge(ctx context.Context, scitizenID string, badgeType string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO badges (id, scitizen_id, badge_type) VALUES ($1, $2, $3)
		 ON CONFLICT (scitizen_id, badge_type) DO NOTHING`,
		ulid.Make().String(), scitizenID, badgeType,
	)
	if err != nil {
		return fmt.Errorf("award badge: %w", err)
	}
	return nil
}

func (r *pgRepo) doGetBadges(ctx context.Context, scitizenID string) ([]Badge, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, scitizen_id, badge_type, awarded_at FROM badges WHERE scitizen_id = $1 ORDER BY awarded_at`,
		scitizenID,
	)
	if err != nil {
		return nil, fmt.Errorf("get badges: %w", err)
	}
	defer rows.Close()

	var badges []Badge
	for rows.Next() {
		var b Badge
		if err := rows.Scan(&b.ID, &b.ScitizenID, &b.BadgeType, &b.AwardedAt); err != nil {
			return nil, fmt.Errorf("scan badge: %w", err)
		}
		badges = append(badges, b)
	}
	return badges, rows.Err()
}

func (r *pgRepo) doGrantSweepstakes(ctx context.Context, input GrantSweepstakesInput) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sweepstakes_entries (id, scitizen_id, entries, milestone_trigger) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (scitizen_id, milestone_trigger) DO NOTHING`,
		ulid.Make().String(), input.ScitizenID, input.Entries, input.MilestoneTrigger,
	)
	if err != nil {
		return fmt.Errorf("grant sweepstakes: %w", err)
	}
	return nil
}

func (r *pgRepo) doGetSweepstakes(ctx context.Context, scitizenID string) ([]SweepstakesEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, scitizen_id, entries, milestone_trigger, granted_at
		 FROM sweepstakes_entries WHERE scitizen_id = $1 ORDER BY granted_at`,
		scitizenID,
	)
	if err != nil {
		return nil, fmt.Errorf("get sweepstakes: %w", err)
	}
	defer rows.Close()

	var entries []SweepstakesEntry
	for rows.Next() {
		var e SweepstakesEntry
		if err := rows.Scan(&e.ID, &e.ScitizenID, &e.Entries, &e.MilestoneTrigger, &e.GrantedAt); err != nil {
			return nil, fmt.Errorf("scan sweepstakes: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
