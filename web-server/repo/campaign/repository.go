package campaign

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type response[T any] struct {
	val T
	err error
}

type createReq struct {
	ctx   context.Context
	input CreateCampaignInput
	resp  chan response[*Campaign]
}

type publishReq struct {
	ctx  context.Context
	id   string
	resp chan response[struct{}]
}

type listReq struct {
	ctx   context.Context
	input ListCampaignsInput
	resp  chan response[[]Campaign]
}

type getRulesReq struct {
	ctx        context.Context
	campaignID string
	resp       chan response[*CampaignRules]
}

type getEligibilityReq struct {
	ctx        context.Context
	campaignID string
	resp       chan response[[]EligibilityCriteria]
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool             *pgxpool.Pool
	createCh         chan createReq
	publishCh        chan publishReq
	listCh           chan listReq
	getRulesCh       chan getRulesReq
	getEligibilityCh chan getEligibilityReq
	shutdownCh       chan shutdownReq
}

// NewRepository creates a campaign repository backed by Postgres.
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:             pool,
		createCh:         make(chan createReq),
		publishCh:        make(chan publishReq),
		listCh:           make(chan listReq),
		getRulesCh:       make(chan getRulesReq),
		getEligibilityCh: make(chan getEligibilityReq),
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
			req.resp <- response[*Campaign]{val: val, err: err}

		case req := <-r.publishCh:
			err := r.doPublish(req.ctx, req.id)
			req.resp <- response[struct{}]{err: err}

		case req := <-r.listCh:
			val, err := r.doList(req.ctx, req.input)
			req.resp <- response[[]Campaign]{val: val, err: err}

		case req := <-r.getRulesCh:
			val, err := r.doGetRules(req.ctx, req.campaignID)
			req.resp <- response[*CampaignRules]{val: val, err: err}

		case req := <-r.getEligibilityCh:
			val, err := r.doGetEligibility(req.ctx, req.campaignID)
			req.resp <- response[[]EligibilityCriteria]{val: val, err: err}

		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) Create(ctx context.Context, input CreateCampaignInput) (*Campaign, error) {
	resp := make(chan response[*Campaign], 1)
	r.createCh <- createReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Publish(ctx context.Context, id string) error {
	resp := make(chan response[struct{}], 1)
	r.publishCh <- publishReq{ctx: ctx, id: id, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) List(ctx context.Context, input ListCampaignsInput) ([]Campaign, error) {
	resp := make(chan response[[]Campaign], 1)
	r.listCh <- listReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetRules(ctx context.Context, campaignID string) (*CampaignRules, error) {
	resp := make(chan response[*CampaignRules], 1)
	r.getRulesCh <- getRulesReq{ctx: ctx, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetEligibility(ctx context.Context, campaignID string) ([]EligibilityCriteria, error) {
	resp := make(chan response[[]EligibilityCriteria], 1)
	r.getEligibilityCh <- getEligibilityReq{ctx: ctx, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doCreate(ctx context.Context, input CreateCampaignInput) (*Campaign, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var c Campaign
	err = tx.QueryRow(ctx,
		`INSERT INTO campaigns (org_id, window_start, window_end, created_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, org_id, status, window_start, window_end, created_by, created_at`,
		input.OrgID, input.WindowStart, input.WindowEnd, input.CreatedBy,
	).Scan(&c.ID, &c.OrgID, &c.Status, &c.WindowStart, &c.WindowEnd, &c.CreatedBy, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert campaign: %w", err)
	}

	for _, p := range input.Parameters {
		_, err := tx.Exec(ctx,
			`INSERT INTO campaign_parameters (campaign_id, name, unit, min_range, max_range, precision)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			c.ID, p.Name, p.Unit, p.MinRange, p.MaxRange, p.Precision,
		)
		if err != nil {
			return nil, fmt.Errorf("insert parameter: %w", err)
		}
	}

	for _, reg := range input.Regions {
		_, err := tx.Exec(ctx,
			`INSERT INTO campaign_regions (campaign_id, geometry) VALUES ($1, $2::jsonb)`,
			c.ID, reg.GeoJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("insert region: %w", err)
		}
	}

	for _, e := range input.Eligibility {
		_, err := tx.Exec(ctx,
			`INSERT INTO campaign_eligibility (campaign_id, device_class, tier, required_sensors, firmware_min)
			 VALUES ($1, $2, $3, $4, $5)`,
			c.ID, e.DeviceClass, e.Tier, e.RequiredSensors, e.FirmwareMin,
		)
		if err != nil {
			return nil, fmt.Errorf("insert eligibility: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &c, nil
}

func (r *pgRepo) doPublish(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE campaigns SET status = 'published' WHERE id = $1 AND status = 'draft'`,
		id,
	)
	if err != nil {
		return fmt.Errorf("publish campaign: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("campaign %s not found or not in draft status", id)
	}
	return nil
}

func (r *pgRepo) doList(ctx context.Context, input ListCampaignsInput) ([]Campaign, error) {
	query := `SELECT id, org_id, status, window_start, window_end, created_by, created_at FROM campaigns WHERE 1=1`
	args := []any{}
	argIdx := 1

	if input.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, input.Status)
		argIdx++
	}
	if input.OrgID != "" {
		query += fmt.Sprintf(" AND org_id = $%d", argIdx)
		args = append(args, input.OrgID)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		var c Campaign
		if err := rows.Scan(&c.ID, &c.OrgID, &c.Status, &c.WindowStart, &c.WindowEnd, &c.CreatedBy, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan campaign: %w", err)
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, rows.Err()
}

func (r *pgRepo) doGetRules(ctx context.Context, campaignID string) (*CampaignRules, error) {
	rules := &CampaignRules{CampaignID: campaignID}

	err := r.pool.QueryRow(ctx,
		`SELECT window_start, window_end FROM campaigns WHERE id = $1`,
		campaignID,
	).Scan(&rules.WindowStart, &rules.WindowEnd)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign %s not found", campaignID)
		}
		return nil, fmt.Errorf("get campaign window: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT name, unit, min_range, max_range, precision FROM campaign_parameters WHERE campaign_id = $1`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get parameters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Parameter
		if err := rows.Scan(&p.Name, &p.Unit, &p.MinRange, &p.MaxRange, &p.Precision); err != nil {
			return nil, fmt.Errorf("scan parameter: %w", err)
		}
		rules.Parameters = append(rules.Parameters, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	regionRows, err := r.pool.Query(ctx,
		`SELECT geometry::text FROM campaign_regions WHERE campaign_id = $1`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get regions: %w", err)
	}
	defer regionRows.Close()

	for regionRows.Next() {
		var reg Region
		if err := regionRows.Scan(&reg.GeoJSON); err != nil {
			return nil, fmt.Errorf("scan region: %w", err)
		}
		rules.Regions = append(rules.Regions, reg)
	}
	return rules, regionRows.Err()
}

func (r *pgRepo) doGetEligibility(ctx context.Context, campaignID string) ([]EligibilityCriteria, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT device_class, tier, required_sensors, firmware_min FROM campaign_eligibility WHERE campaign_id = $1`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get eligibility: %w", err)
	}
	defer rows.Close()

	var criteria []EligibilityCriteria
	for rows.Next() {
		var e EligibilityCriteria
		if err := rows.Scan(&e.DeviceClass, &e.Tier, &e.RequiredSensors, &e.FirmwareMin); err != nil {
			return nil, fmt.Errorf("scan eligibility: %w", err)
		}
		criteria = append(criteria, e)
	}
	return criteria, rows.Err()
}

