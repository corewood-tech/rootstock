package enrollment

import (
	"context"
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

type enrollReq struct {
	ctx   context.Context
	input EnrollInput
	resp  chan response[*Enrollment]
}

type withdrawReq struct {
	ctx  context.Context
	id   string
	resp chan response[struct{}]
}

type getByIDReq struct {
	ctx  context.Context
	id   string
	resp chan response[*Enrollment]
}

type getByDeviceCampaignReq struct {
	ctx        context.Context
	deviceID   string
	campaignID string
	resp       chan response[*Enrollment]
}

type markReadReq struct {
	ctx    context.Context
	userID string
	ids    []string
	resp   chan response[int]
}

type createNotificationReq struct {
	ctx   context.Context
	input CreateNotificationInput
	resp  chan response[struct{}]
}

type getPreferencesReq struct {
	ctx    context.Context
	userID string
	resp   chan response[[]NotificationPreference]
}

type updatePreferencesReq struct {
	ctx    context.Context
	userID string
	prefs  []NotificationPreference
	resp   chan response[struct{}]
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool                  *pgxpool.Pool
	enrollCh              chan enrollReq
	withdrawCh            chan withdrawReq
	getByIDCh             chan getByIDReq
	getByDeviceCampaignCh chan getByDeviceCampaignReq
	markReadCh            chan markReadReq
	createNotificationCh  chan createNotificationReq
	getPreferencesCh      chan getPreferencesReq
	updatePreferencesCh   chan updatePreferencesReq
	shutdownCh            chan shutdownReq
}

// NewRepository creates an enrollment repository backed by Postgres.
// Graph node: 0x2c (EnrollmentRepository)
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:                  pool,
		enrollCh:              make(chan enrollReq),
		withdrawCh:            make(chan withdrawReq),
		getByIDCh:             make(chan getByIDReq),
		getByDeviceCampaignCh: make(chan getByDeviceCampaignReq),
		markReadCh:            make(chan markReadReq),
		createNotificationCh:  make(chan createNotificationReq),
		getPreferencesCh:      make(chan getPreferencesReq),
		updatePreferencesCh:   make(chan updatePreferencesReq),
		shutdownCh:            make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *pgRepo) manage() {
	for {
		select {
		case req := <-r.enrollCh:
			val, err := r.doEnroll(req.ctx, req.input)
			req.resp <- response[*Enrollment]{val: val, err: err}
		case req := <-r.withdrawCh:
			err := r.doWithdraw(req.ctx, req.id)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.getByIDCh:
			val, err := r.doGetByID(req.ctx, req.id)
			req.resp <- response[*Enrollment]{val: val, err: err}
		case req := <-r.getByDeviceCampaignCh:
			val, err := r.doGetByDeviceCampaign(req.ctx, req.deviceID, req.campaignID)
			req.resp <- response[*Enrollment]{val: val, err: err}
		case req := <-r.markReadCh:
			val, err := r.doMarkRead(req.ctx, req.userID, req.ids)
			req.resp <- response[int]{val: val, err: err}
		case req := <-r.createNotificationCh:
			err := r.doCreateNotification(req.ctx, req.input)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.getPreferencesCh:
			val, err := r.doGetPreferences(req.ctx, req.userID)
			req.resp <- response[[]NotificationPreference]{val: val, err: err}
		case req := <-r.updatePreferencesCh:
			err := r.doUpdatePreferences(req.ctx, req.userID, req.prefs)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) Enroll(ctx context.Context, input EnrollInput) (*Enrollment, error) {
	resp := make(chan response[*Enrollment], 1)
	r.enrollCh <- enrollReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Withdraw(ctx context.Context, enrollmentID string) error {
	resp := make(chan response[struct{}], 1)
	r.withdrawCh <- withdrawReq{ctx: ctx, id: enrollmentID, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) GetByID(ctx context.Context, id string) (*Enrollment, error) {
	resp := make(chan response[*Enrollment], 1)
	r.getByIDCh <- getByIDReq{ctx: ctx, id: id, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetByDeviceCampaign(ctx context.Context, deviceID, campaignID string) (*Enrollment, error) {
	resp := make(chan response[*Enrollment], 1)
	r.getByDeviceCampaignCh <- getByDeviceCampaignReq{ctx: ctx, deviceID: deviceID, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) MarkRead(ctx context.Context, userID string, ids []string) (int, error) {
	resp := make(chan response[int], 1)
	r.markReadCh <- markReadReq{ctx: ctx, userID: userID, ids: ids, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) CreateNotification(ctx context.Context, input CreateNotificationInput) error {
	resp := make(chan response[struct{}], 1)
	r.createNotificationCh <- createNotificationReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) GetPreferences(ctx context.Context, userID string) ([]NotificationPreference, error) {
	resp := make(chan response[[]NotificationPreference], 1)
	r.getPreferencesCh <- getPreferencesReq{ctx: ctx, userID: userID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) UpdatePreferences(ctx context.Context, userID string, prefs []NotificationPreference) error {
	resp := make(chan response[struct{}], 1)
	r.updatePreferencesCh <- updatePreferencesReq{ctx: ctx, userID: userID, prefs: prefs, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doEnroll(ctx context.Context, input EnrollInput) (*Enrollment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	enrollmentID := ulid.Make().String()
	var e Enrollment
	err = tx.QueryRow(ctx,
		`INSERT INTO campaign_enrollments (id, device_id, campaign_id, scitizen_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, device_id, campaign_id, scitizen_id, status, enrolled_at`,
		enrollmentID, input.DeviceID, input.CampaignID, input.ScitizenID,
	).Scan(&e.ID, &e.DeviceID, &e.CampaignID, &e.ScitizenID, &e.Status, &e.EnrolledAt)
	if err != nil {
		return nil, fmt.Errorf("insert enrollment: %w", err)
	}

	// Record consent
	consentID := ulid.Make().String()
	_, err = tx.Exec(ctx,
		`INSERT INTO consent_records (id, enrollment_id, version, scope)
		 VALUES ($1, $2, $3, $4)`,
		consentID, enrollmentID, input.ConsentVersion, input.ConsentScope,
	)
	if err != nil {
		return nil, fmt.Errorf("insert consent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &e, nil
}

func (r *pgRepo) doWithdraw(ctx context.Context, enrollmentID string) error {
	now := time.Now()
	tag, err := r.pool.Exec(ctx,
		`UPDATE campaign_enrollments SET status = 'withdrawn', withdrawn_at = $1
		 WHERE id = $2 AND status = 'active'`,
		now, enrollmentID,
	)
	if err != nil {
		return fmt.Errorf("withdraw enrollment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("enrollment %s not found or already withdrawn", enrollmentID)
	}
	return nil
}

func (r *pgRepo) doGetByID(ctx context.Context, id string) (*Enrollment, error) {
	var e Enrollment
	err := r.pool.QueryRow(ctx,
		`SELECT id, device_id, campaign_id, scitizen_id, status, enrolled_at, withdrawn_at
		 FROM campaign_enrollments WHERE id = $1`,
		id,
	).Scan(&e.ID, &e.DeviceID, &e.CampaignID, &e.ScitizenID, &e.Status, &e.EnrolledAt, &e.WithdrawnAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("enrollment %s not found", id)
		}
		return nil, fmt.Errorf("get enrollment: %w", err)
	}
	return &e, nil
}

func (r *pgRepo) doGetByDeviceCampaign(ctx context.Context, deviceID, campaignID string) (*Enrollment, error) {
	var e Enrollment
	err := r.pool.QueryRow(ctx,
		`SELECT id, device_id, campaign_id, scitizen_id, status, enrolled_at, withdrawn_at
		 FROM campaign_enrollments WHERE device_id = $1 AND campaign_id = $2`,
		deviceID, campaignID,
	).Scan(&e.ID, &e.DeviceID, &e.CampaignID, &e.ScitizenID, &e.Status, &e.EnrolledAt, &e.WithdrawnAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get enrollment by device+campaign: %w", err)
	}
	return &e, nil
}

func (r *pgRepo) doMarkRead(ctx context.Context, userID string, ids []string) (int, error) {
	tag, err := r.pool.Exec(ctx,
		`UPDATE notifications SET read = true
		 WHERE user_id = $1 AND id = ANY($2) AND read = false`,
		userID, ids,
	)
	if err != nil {
		return 0, fmt.Errorf("mark read: %w", err)
	}
	return int(tag.RowsAffected()), nil
}

func (r *pgRepo) doCreateNotification(ctx context.Context, input CreateNotificationInput) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications (id, user_id, type, message, resource_link)
		 VALUES ($1, $2, $3, $4, $5)`,
		ulid.Make().String(), input.UserID, input.Type, input.Message, input.ResourceLink,
	)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	return nil
}

func (r *pgRepo) doGetPreferences(ctx context.Context, userID string) ([]NotificationPreference, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT type, in_app, email FROM notification_preferences WHERE user_id = $1 ORDER BY type`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	defer rows.Close()

	var prefs []NotificationPreference
	for rows.Next() {
		var p NotificationPreference
		if err := rows.Scan(&p.Type, &p.InApp, &p.Email); err != nil {
			return nil, fmt.Errorf("scan preference: %w", err)
		}
		prefs = append(prefs, p)
	}
	return prefs, rows.Err()
}

func (r *pgRepo) doUpdatePreferences(ctx context.Context, userID string, prefs []NotificationPreference) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, p := range prefs {
		_, err := tx.Exec(ctx,
			`INSERT INTO notification_preferences (user_id, type, in_app, email)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (user_id, type) DO UPDATE SET in_app = $3, email = $4`,
			userID, p.Type, p.InApp, p.Email,
		)
		if err != nil {
			return fmt.Errorf("upsert preference: %w", err)
		}
	}

	return tx.Commit(ctx)
}
