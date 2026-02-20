package reading

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

type persistReq struct {
	ctx   context.Context
	input PersistReadingInput
	resp  chan response[*Reading]
}

type quarantineReq struct {
	ctx    context.Context
	id     string
	reason string
	resp   chan response[struct{}]
}

type queryReq struct {
	ctx   context.Context
	input QueryReadingsInput
	resp  chan response[[]Reading]
}

type quarantineByWindowReq struct {
	ctx   context.Context
	input QuarantineByWindowInput
	resp  chan response[int64]
}

type getQualityReq struct {
	ctx        context.Context
	campaignID string
	resp       chan response[*QualityMetrics]
}

type getScitizenStatsReq struct {
	ctx        context.Context
	scitizenID string
	resp       chan response[*ScitizenReadingStats]
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool                  *pgxpool.Pool
	persistCh             chan persistReq
	quarantineCh          chan quarantineReq
	queryCh               chan queryReq
	quarantineByWindowCh  chan quarantineByWindowReq
	getQualityCh          chan getQualityReq
	getScitizenStatsCh    chan getScitizenStatsReq
	shutdownCh            chan shutdownReq
}

// NewRepository creates a reading repository backed by Postgres.
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:                 pool,
		persistCh:            make(chan persistReq),
		quarantineCh:         make(chan quarantineReq),
		queryCh:              make(chan queryReq),
		quarantineByWindowCh: make(chan quarantineByWindowReq),
		getQualityCh:         make(chan getQualityReq),
		getScitizenStatsCh:   make(chan getScitizenStatsReq),
		shutdownCh:           make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *pgRepo) manage() {
	for {
		select {
		case req := <-r.persistCh:
			val, err := r.doPersist(req.ctx, req.input)
			req.resp <- response[*Reading]{val: val, err: err}
		case req := <-r.quarantineCh:
			err := r.doQuarantine(req.ctx, req.id, req.reason)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.queryCh:
			val, err := r.doQuery(req.ctx, req.input)
			req.resp <- response[[]Reading]{val: val, err: err}
		case req := <-r.quarantineByWindowCh:
			val, err := r.doQuarantineByWindow(req.ctx, req.input)
			req.resp <- response[int64]{val: val, err: err}
		case req := <-r.getQualityCh:
			val, err := r.doGetQuality(req.ctx, req.campaignID)
			req.resp <- response[*QualityMetrics]{val: val, err: err}
		case req := <-r.getScitizenStatsCh:
			val, err := r.doGetScitizenReadingStats(req.ctx, req.scitizenID)
			req.resp <- response[*ScitizenReadingStats]{val: val, err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) Persist(ctx context.Context, input PersistReadingInput) (*Reading, error) {
	resp := make(chan response[*Reading], 1)
	r.persistCh <- persistReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Quarantine(ctx context.Context, id string, reason string) error {
	resp := make(chan response[struct{}], 1)
	r.quarantineCh <- quarantineReq{ctx: ctx, id: id, reason: reason, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) Query(ctx context.Context, input QueryReadingsInput) ([]Reading, error) {
	resp := make(chan response[[]Reading], 1)
	r.queryCh <- queryReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) QuarantineByWindow(ctx context.Context, input QuarantineByWindowInput) (int64, error) {
	resp := make(chan response[int64], 1)
	r.quarantineByWindowCh <- quarantineByWindowReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetCampaignQuality(ctx context.Context, campaignID string) (*QualityMetrics, error) {
	resp := make(chan response[*QualityMetrics], 1)
	r.getQualityCh <- getQualityReq{ctx: ctx, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetScitizenReadingStats(ctx context.Context, scitizenID string) (*ScitizenReadingStats, error) {
	resp := make(chan response[*ScitizenReadingStats], 1)
	r.getScitizenStatsCh <- getScitizenStatsReq{ctx: ctx, scitizenID: scitizenID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doPersist(ctx context.Context, input PersistReadingInput) (*Reading, error) {
	var geo *string
	if input.Geolocation != "" {
		geo = &input.Geolocation
	}

	var rd Reading
	err := r.pool.QueryRow(ctx,
		`INSERT INTO readings (id, device_id, campaign_id, value, timestamp, geolocation, firmware_version, cert_serial)
		 VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8)
		 RETURNING id, device_id, campaign_id, value, timestamp, geolocation::text, firmware_version, cert_serial, ingested_at, status, quarantine_reason`,
		ulid.Make().String(), input.DeviceID, input.CampaignID, input.Value, input.Timestamp, geo, input.FirmwareVersion, input.CertSerial,
	).Scan(&rd.ID, &rd.DeviceID, &rd.CampaignID, &rd.Value, &rd.Timestamp, &rd.Geolocation, &rd.FirmwareVersion, &rd.CertSerial, &rd.IngestedAt, &rd.Status, &rd.QuarantineReason)
	if err != nil {
		return nil, fmt.Errorf("insert reading: %w", err)
	}
	return &rd, nil
}

func (r *pgRepo) doQuarantine(ctx context.Context, id string, reason string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE readings SET status = 'quarantined', quarantine_reason = $1 WHERE id = $2`,
		reason, id,
	)
	if err != nil {
		return fmt.Errorf("quarantine reading: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reading %s not found", id)
	}
	return nil
}

func (r *pgRepo) doQuery(ctx context.Context, input QueryReadingsInput) ([]Reading, error) {
	query := `SELECT id, device_id, campaign_id, value, timestamp, geolocation::text, firmware_version, cert_serial, ingested_at, status, quarantine_reason
	          FROM readings WHERE 1=1`
	args := []any{}
	argIdx := 1

	if input.CampaignID != "" {
		query += fmt.Sprintf(" AND campaign_id = $%d", argIdx)
		args = append(args, input.CampaignID)
		argIdx++
	}
	if input.DeviceID != "" {
		query += fmt.Sprintf(" AND device_id = $%d", argIdx)
		args = append(args, input.DeviceID)
		argIdx++
	}
	if input.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, input.Status)
		argIdx++
	}
	if input.Since != nil {
		query += fmt.Sprintf(" AND timestamp >= $%d", argIdx)
		args = append(args, input.Since)
		argIdx++
	}
	if input.Until != nil {
		query += fmt.Sprintf(" AND timestamp <= $%d", argIdx)
		args = append(args, input.Until)
		argIdx++
	}

	query += " ORDER BY timestamp DESC"

	if input.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, input.Limit)
		argIdx++
	}

	if input.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, input.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query readings: %w", err)
	}
	defer rows.Close()

	var readings []Reading
	for rows.Next() {
		var rd Reading
		if err := rows.Scan(&rd.ID, &rd.DeviceID, &rd.CampaignID, &rd.Value, &rd.Timestamp, &rd.Geolocation, &rd.FirmwareVersion, &rd.CertSerial, &rd.IngestedAt, &rd.Status, &rd.QuarantineReason); err != nil {
			return nil, fmt.Errorf("scan reading: %w", err)
		}
		readings = append(readings, rd)
	}
	return readings, rows.Err()
}

func (r *pgRepo) doQuarantineByWindow(ctx context.Context, input QuarantineByWindowInput) (int64, error) {
	if len(input.DeviceIDs) == 0 {
		return 0, nil
	}

	// Build ANY($N) for device IDs
	tag, err := r.pool.Exec(ctx,
		`UPDATE readings
		 SET status = 'quarantined', quarantine_reason = $1
		 WHERE device_id = ANY($2)
		   AND timestamp >= $3
		   AND timestamp <= $4
		   AND status = 'accepted'`,
		input.Reason, input.DeviceIDs, input.Since, input.Until,
	)
	if err != nil {
		return 0, fmt.Errorf("quarantine by window: %w", err)
	}
	return tag.RowsAffected(), nil
}

func (r *pgRepo) doGetQuality(ctx context.Context, campaignID string) (*QualityMetrics, error) {
	q := &QualityMetrics{CampaignID: campaignID}
	err := r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'quarantined' THEN 1 ELSE 0 END), 0)
		 FROM readings WHERE campaign_id = $1`,
		campaignID,
	).Scan(&q.AcceptedCount, &q.QuarantineCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return q, nil
		}
		return nil, fmt.Errorf("get quality: %w", err)
	}
	return q, nil
}

func (r *pgRepo) doGetScitizenReadingStats(ctx context.Context, scitizenID string) (*ScitizenReadingStats, error) {
	s := &ScitizenReadingStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(COUNT(*) FILTER (WHERE r.status = 'accepted'), 0),
			COALESCE(COUNT(*) FILTER (WHERE r.status = 'accepted')::float / NULLIF(COUNT(*), 0), 0),
			COALESCE(
				COUNT(DISTINCT DATE(r.timestamp))::float /
				GREATEST(EXTRACT(EPOCH FROM (now() - MIN(r.ingested_at))) / 86400, 1),
			0),
			COALESCE(COUNT(DISTINCT r.campaign_id), 0)
		 FROM readings r
		 JOIN devices d ON d.id = r.device_id
		 WHERE d.owner_id = $1`,
		scitizenID,
	).Scan(&s.Volume, &s.QualityRate, &s.Consistency, &s.Diversity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return s, nil
		}
		return nil, fmt.Errorf("get scitizen reading stats: %w", err)
	}
	return s, nil
}
