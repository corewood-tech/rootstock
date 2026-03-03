package reading

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

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

type getDeviceBreakdownReq struct {
	ctx        context.Context
	campaignID string
	hmacSecret string
	resp       chan response[[]DeviceBreakdown]
}

type getTemporalCoverageReq struct {
	ctx        context.Context
	campaignID string
	resp       chan response[[]TemporalBucket]
}

type getEnrollmentFunnelReq struct {
	ctx        context.Context
	campaignID string
	resp       chan response[*EnrollmentFunnel]
}

type getScitizenStatsReq struct {
	ctx        context.Context
	scitizenID string
	resp       chan response[*ScitizenReadingStats]
}

type shutdownReq struct {
	resp chan struct{}
}

type quarantineValueReq struct {
	ctx    context.Context
	id     string
	reason string
	resp   chan response[struct{}]
}

type pgRepo struct {
	pool                    *pgxpool.Pool
	persistCh               chan persistReq
	quarantineCh            chan quarantineReq
	quarantineValueCh       chan quarantineValueReq
	queryCh                 chan queryReq
	quarantineByWindowCh    chan quarantineByWindowReq
	getQualityCh            chan getQualityReq
	getDeviceBreakdownCh    chan getDeviceBreakdownReq
	getTemporalCoverageCh   chan getTemporalCoverageReq
	getEnrollmentFunnelCh   chan getEnrollmentFunnelReq
	getScitizenStatsCh      chan getScitizenStatsReq
	shutdownCh              chan shutdownReq
}

// NewRepository creates a reading repository backed by Postgres.
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:                    pool,
		persistCh:               make(chan persistReq),
		quarantineCh:            make(chan quarantineReq),
		quarantineValueCh:       make(chan quarantineValueReq),
		queryCh:                 make(chan queryReq),
		quarantineByWindowCh:    make(chan quarantineByWindowReq),
		getQualityCh:            make(chan getQualityReq),
		getDeviceBreakdownCh:    make(chan getDeviceBreakdownReq),
		getTemporalCoverageCh:   make(chan getTemporalCoverageReq),
		getEnrollmentFunnelCh:   make(chan getEnrollmentFunnelReq),
		getScitizenStatsCh:      make(chan getScitizenStatsReq),
		shutdownCh:              make(chan shutdownReq),
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
		case req := <-r.quarantineValueCh:
			err := r.doQuarantineValue(req.ctx, req.id, req.reason)
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
		case req := <-r.getDeviceBreakdownCh:
			val, err := r.doGetDeviceBreakdown(req.ctx, req.campaignID, req.hmacSecret)
			req.resp <- response[[]DeviceBreakdown]{val: val, err: err}
		case req := <-r.getTemporalCoverageCh:
			val, err := r.doGetTemporalCoverage(req.ctx, req.campaignID)
			req.resp <- response[[]TemporalBucket]{val: val, err: err}
		case req := <-r.getEnrollmentFunnelCh:
			val, err := r.doGetEnrollmentFunnel(req.ctx, req.campaignID)
			req.resp <- response[*EnrollmentFunnel]{val: val, err: err}
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

func (r *pgRepo) QuarantineValue(ctx context.Context, readingValueID string, reason string) error {
	resp := make(chan response[struct{}], 1)
	r.quarantineValueCh <- quarantineValueReq{ctx: ctx, id: readingValueID, reason: reason, resp: resp}
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

func (r *pgRepo) GetCampaignDeviceBreakdown(ctx context.Context, campaignID string, hmacSecret string) ([]DeviceBreakdown, error) {
	resp := make(chan response[[]DeviceBreakdown], 1)
	r.getDeviceBreakdownCh <- getDeviceBreakdownReq{ctx: ctx, campaignID: campaignID, hmacSecret: hmacSecret, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetCampaignTemporalCoverage(ctx context.Context, campaignID string) ([]TemporalBucket, error) {
	resp := make(chan response[[]TemporalBucket], 1)
	r.getTemporalCoverageCh <- getTemporalCoverageReq{ctx: ctx, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetEnrollmentFunnel(ctx context.Context, campaignID string) (*EnrollmentFunnel, error) {
	resp := make(chan response[*EnrollmentFunnel], 1)
	r.getEnrollmentFunnelCh <- getEnrollmentFunnelReq{ctx: ctx, campaignID: campaignID, resp: resp}
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

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var rd Reading
	readingID := ulid.Make().String()
	err = tx.QueryRow(ctx,
		`INSERT INTO readings (id, device_id, campaign_id, value, timestamp, geolocation, firmware_version, cert_serial)
		 VALUES ($1, $2, $3, NULL, $4, $5::jsonb, $6, $7)
		 RETURNING id, device_id, campaign_id, value, timestamp, geolocation::text, firmware_version, cert_serial, ingested_at, status, quarantine_reason`,
		readingID, input.DeviceID, input.CampaignID, input.Timestamp, geo, input.FirmwareVersion, input.CertSerial,
	).Scan(&rd.ID, &rd.DeviceID, &rd.CampaignID, &rd.Value, &rd.Timestamp, &rd.Geolocation, &rd.FirmwareVersion, &rd.CertSerial, &rd.IngestedAt, &rd.Status, &rd.QuarantineReason)
	if err != nil {
		return nil, fmt.Errorf("insert reading: %w", err)
	}

	for _, v := range input.Values {
		var rv ReadingValue
		err = tx.QueryRow(ctx,
			`INSERT INTO reading_values (id, reading_id, parameter_name, value)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id, reading_id, parameter_name, value, status, quarantine_reason`,
			ulid.Make().String(), readingID, v.ParameterName, v.Value,
		).Scan(&rv.ID, &rv.ReadingID, &rv.ParameterName, &rv.Value, &rv.Status, &rv.QuarantineReason)
		if err != nil {
			return nil, fmt.Errorf("insert reading value %s: %w", v.ParameterName, err)
		}
		rd.Values = append(rd.Values, rv)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load reading_values for each reading
	for i := range readings {
		valRows, err := r.pool.Query(ctx,
			`SELECT id, reading_id, parameter_name, value, status, quarantine_reason
			 FROM reading_values WHERE reading_id = $1`,
			readings[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf("query reading values: %w", err)
		}
		for valRows.Next() {
			var rv ReadingValue
			if err := valRows.Scan(&rv.ID, &rv.ReadingID, &rv.ParameterName, &rv.Value, &rv.Status, &rv.QuarantineReason); err != nil {
				valRows.Close()
				return nil, fmt.Errorf("scan reading value: %w", err)
			}
			readings[i].Values = append(readings[i].Values, rv)
		}
		valRows.Close()
		if err := valRows.Err(); err != nil {
			return nil, err
		}
	}

	return readings, nil
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

func (r *pgRepo) doQuarantineValue(ctx context.Context, readingValueID string, reason string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE reading_values SET status = 'quarantined', quarantine_reason = $1 WHERE id = $2`,
		reason, readingValueID,
	)
	if err != nil {
		return fmt.Errorf("quarantine reading value: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reading value %s not found", readingValueID)
	}
	return nil
}

func (r *pgRepo) doGetQuality(ctx context.Context, campaignID string) (*QualityMetrics, error) {
	q := &QualityMetrics{CampaignID: campaignID}

	// Overall counts from reading_values
	err := r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(SUM(CASE WHEN rv.status = 'accepted' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN rv.status = 'quarantined' THEN 1 ELSE 0 END), 0)
		 FROM reading_values rv
		 JOIN readings r ON r.id = rv.reading_id
		 WHERE r.campaign_id = $1`,
		campaignID,
	).Scan(&q.AcceptedCount, &q.QuarantineCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return q, nil
		}
		return nil, fmt.Errorf("get quality: %w", err)
	}

	// Per-parameter breakdown
	rows, err := r.pool.Query(ctx,
		`SELECT
			rv.parameter_name,
			COALESCE(SUM(CASE WHEN rv.status = 'accepted' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN rv.status = 'quarantined' THEN 1 ELSE 0 END), 0)
		 FROM reading_values rv
		 JOIN readings r ON r.id = rv.reading_id
		 WHERE r.campaign_id = $1
		 GROUP BY rv.parameter_name
		 ORDER BY rv.parameter_name`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get per-parameter quality: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pq ParameterQuality
		if err := rows.Scan(&pq.ParameterName, &pq.AcceptedCount, &pq.QuarantinedCount); err != nil {
			return nil, fmt.Errorf("scan parameter quality: %w", err)
		}
		q.PerParameter = append(q.PerParameter, pq)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return q, nil
}

func pseudonymizeDeviceID(deviceID, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(deviceID))
	return hex.EncodeToString(mac.Sum(nil))
}

func (r *pgRepo) doGetDeviceBreakdown(ctx context.Context, campaignID string, hmacSecret string) ([]DeviceBreakdown, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT
			r.device_id,
			d.class,
			COALESCE(COUNT(rv.*) FILTER (WHERE rv.status = 'accepted')::float / NULLIF(COUNT(rv.*), 0), 0) AS acceptance_rate,
			COUNT(DISTINCT r.id) AS reading_count,
			MAX(r.timestamp) AS last_seen
		 FROM readings r
		 JOIN devices d ON d.id = r.device_id
		 LEFT JOIN reading_values rv ON rv.reading_id = r.id
		 WHERE r.campaign_id = $1
		 GROUP BY r.device_id, d.class
		 ORDER BY reading_count DESC`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get device breakdown: %w", err)
	}
	defer rows.Close()

	var result []DeviceBreakdown
	for rows.Next() {
		var deviceID, class string
		var acceptanceRate float64
		var readingCount int
		var lastSeen time.Time
		if err := rows.Scan(&deviceID, &class, &acceptanceRate, &readingCount, &lastSeen); err != nil {
			return nil, fmt.Errorf("scan device breakdown: %w", err)
		}
		ls := lastSeen.Format(time.RFC3339)
		result = append(result, DeviceBreakdown{
			PseudoDeviceID: pseudonymizeDeviceID(deviceID, hmacSecret),
			DeviceClass:    class,
			AcceptanceRate: acceptanceRate,
			ReadingCount:   readingCount,
			LastSeen:       &ls,
		})
	}
	return result, rows.Err()
}

func (r *pgRepo) doGetTemporalCoverage(ctx context.Context, campaignID string) ([]TemporalBucket, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT
			date_trunc('hour', r.timestamp) AS bucket,
			COUNT(*) AS cnt
		 FROM readings r
		 WHERE r.campaign_id = $1
		 GROUP BY bucket
		 ORDER BY bucket`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get temporal coverage: %w", err)
	}
	defer rows.Close()

	var result []TemporalBucket
	for rows.Next() {
		var bucket time.Time
		var count int
		if err := rows.Scan(&bucket, &count); err != nil {
			return nil, fmt.Errorf("scan temporal bucket: %w", err)
		}
		result = append(result, TemporalBucket{
			Bucket: bucket.Format(time.RFC3339),
			Count:  count,
		})
	}
	return result, rows.Err()
}

func (r *pgRepo) doGetEnrollmentFunnel(ctx context.Context, campaignID string) (*EnrollmentFunnel, error) {
	f := &EnrollmentFunnel{}
	err := r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(COUNT(*), 0),
			COALESCE(COUNT(*) FILTER (WHERE ce.status = 'active'), 0),
			COALESCE(COUNT(DISTINCT ce.device_id) FILTER (
				WHERE ce.status = 'active' AND EXISTS (
					SELECT 1 FROM readings r WHERE r.device_id = ce.device_id AND r.campaign_id = ce.campaign_id
				)
			), 0)
		 FROM campaign_enrollments ce
		 WHERE ce.campaign_id = $1`,
		campaignID,
	).Scan(&f.Enrolled, &f.Active, &f.Contributing)
	if err != nil {
		if err == pgx.ErrNoRows {
			return f, nil
		}
		return nil, fmt.Errorf("get enrollment funnel: %w", err)
	}
	return f, nil
}

func (r *pgRepo) doGetScitizenReadingStats(ctx context.Context, scitizenID string) (*ScitizenReadingStats, error) {
	s := &ScitizenReadingStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
			COALESCE(COUNT(rv.*) FILTER (WHERE rv.status = 'accepted'), 0),
			COALESCE(COUNT(rv.*) FILTER (WHERE rv.status = 'accepted')::float / NULLIF(COUNT(rv.*), 0), 0),
			COALESCE(
				COUNT(DISTINCT DATE(r.timestamp))::float /
				GREATEST(EXTRACT(EPOCH FROM (now() - MIN(r.ingested_at))) / 86400, 1),
			0),
			COALESCE(COUNT(DISTINCT r.campaign_id), 0)
		 FROM readings r
		 JOIN devices d ON d.id = r.device_id
		 LEFT JOIN reading_values rv ON rv.reading_id = r.id
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
