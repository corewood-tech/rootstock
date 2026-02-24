package scitizen

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

type createProfileReq struct {
	ctx   context.Context
	input CreateProfileInput
	resp  chan response[*Profile]
}

type getProfileReq struct {
	ctx    context.Context
	userID string
	resp   chan response[*Profile]
}

type updateOnboardingReq struct {
	ctx   context.Context
	input UpdateOnboardingInput
	resp  chan response[struct{}]
}

type getDashboardReq struct {
	ctx    context.Context
	userID string
	resp   chan response[*Dashboard]
}

type getContributionsReq struct {
	ctx    context.Context
	userID string
	resp   chan response[[]ReadingHistory]
}

type browseCampaignsReq struct {
	ctx   context.Context
	input BrowseInput
	resp  chan response[browseResult]
}

type browseResult struct {
	campaigns []CampaignSummary
	total     int
}

type getCampaignDetailReq struct {
	ctx        context.Context
	campaignID string
	resp       chan response[*CampaignDetail]
}

type searchCampaignsReq struct {
	ctx   context.Context
	input SearchInput
	resp  chan response[browseResult]
}

type getDevicesReq struct {
	ctx     context.Context
	ownerID string
	resp    chan response[[]DeviceSummary]
}

type getDeviceDetailReq struct {
	ctx      context.Context
	deviceID string
	resp     chan response[*DeviceDetail]
}

type getNotificationsReq struct {
	ctx   context.Context
	input GetNotificationsInput
	resp  chan response[notificationsResult]
}

type notificationsResult struct {
	notifications []Notification
	unreadCount   int
	total         int
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool                *pgxpool.Pool
	createProfileCh     chan createProfileReq
	getProfileCh        chan getProfileReq
	updateOnboardingCh  chan updateOnboardingReq
	getDashboardCh      chan getDashboardReq
	getContributionsCh  chan getContributionsReq
	browseCampaignsCh   chan browseCampaignsReq
	getCampaignDetailCh chan getCampaignDetailReq
	searchCampaignsCh   chan searchCampaignsReq
	getDevicesCh        chan getDevicesReq
	getDeviceDetailCh   chan getDeviceDetailReq
	getNotificationsCh  chan getNotificationsReq
	shutdownCh          chan shutdownReq
}

// NewRepository creates a scitizen repository backed by Postgres.
// Graph node: 0x24 (ScitizenRepository)
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:                pool,
		createProfileCh:     make(chan createProfileReq),
		getProfileCh:        make(chan getProfileReq),
		updateOnboardingCh:  make(chan updateOnboardingReq),
		getDashboardCh:      make(chan getDashboardReq),
		getContributionsCh:  make(chan getContributionsReq),
		browseCampaignsCh:   make(chan browseCampaignsReq),
		getCampaignDetailCh: make(chan getCampaignDetailReq),
		searchCampaignsCh:   make(chan searchCampaignsReq),
		getDevicesCh:        make(chan getDevicesReq),
		getDeviceDetailCh:   make(chan getDeviceDetailReq),
		getNotificationsCh:  make(chan getNotificationsReq),
		shutdownCh:          make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *pgRepo) manage() {
	for {
		select {
		case req := <-r.createProfileCh:
			val, err := r.doCreateProfile(req.ctx, req.input)
			req.resp <- response[*Profile]{val: val, err: err}
		case req := <-r.getProfileCh:
			val, err := r.doGetProfile(req.ctx, req.userID)
			req.resp <- response[*Profile]{val: val, err: err}
		case req := <-r.updateOnboardingCh:
			err := r.doUpdateOnboarding(req.ctx, req.input)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.getDashboardCh:
			val, err := r.doGetDashboard(req.ctx, req.userID)
			req.resp <- response[*Dashboard]{val: val, err: err}
		case req := <-r.getContributionsCh:
			val, err := r.doGetContributions(req.ctx, req.userID)
			req.resp <- response[[]ReadingHistory]{val: val, err: err}
		case req := <-r.browseCampaignsCh:
			val, err := r.doBrowseCampaigns(req.ctx, req.input)
			req.resp <- response[browseResult]{val: val, err: err}
		case req := <-r.getCampaignDetailCh:
			val, err := r.doGetCampaignDetail(req.ctx, req.campaignID)
			req.resp <- response[*CampaignDetail]{val: val, err: err}
		case req := <-r.searchCampaignsCh:
			val, err := r.doSearchCampaigns(req.ctx, req.input)
			req.resp <- response[browseResult]{val: val, err: err}
		case req := <-r.getDevicesCh:
			val, err := r.doGetDevices(req.ctx, req.ownerID)
			req.resp <- response[[]DeviceSummary]{val: val, err: err}
		case req := <-r.getDeviceDetailCh:
			val, err := r.doGetDeviceDetail(req.ctx, req.deviceID)
			req.resp <- response[*DeviceDetail]{val: val, err: err}
		case req := <-r.getNotificationsCh:
			val, err := r.doGetNotifications(req.ctx, req.input)
			req.resp <- response[notificationsResult]{val: val, err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) CreateProfile(ctx context.Context, input CreateProfileInput) (*Profile, error) {
	resp := make(chan response[*Profile], 1)
	r.createProfileCh <- createProfileReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	resp := make(chan response[*Profile], 1)
	r.getProfileCh <- getProfileReq{ctx: ctx, userID: userID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) UpdateOnboarding(ctx context.Context, input UpdateOnboardingInput) error {
	resp := make(chan response[struct{}], 1)
	r.updateOnboardingCh <- updateOnboardingReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) GetDashboard(ctx context.Context, userID string) (*Dashboard, error) {
	resp := make(chan response[*Dashboard], 1)
	r.getDashboardCh <- getDashboardReq{ctx: ctx, userID: userID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetContributions(ctx context.Context, userID string) ([]ReadingHistory, error) {
	resp := make(chan response[[]ReadingHistory], 1)
	r.getContributionsCh <- getContributionsReq{ctx: ctx, userID: userID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) BrowseCampaigns(ctx context.Context, input BrowseInput) ([]CampaignSummary, int, error) {
	resp := make(chan response[browseResult], 1)
	r.browseCampaignsCh <- browseCampaignsReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val.campaigns, res.val.total, res.err
}

func (r *pgRepo) GetCampaignDetail(ctx context.Context, campaignID string) (*CampaignDetail, error) {
	resp := make(chan response[*CampaignDetail], 1)
	r.getCampaignDetailCh <- getCampaignDetailReq{ctx: ctx, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) SearchCampaigns(ctx context.Context, input SearchInput) ([]CampaignSummary, int, error) {
	resp := make(chan response[browseResult], 1)
	r.searchCampaignsCh <- searchCampaignsReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val.campaigns, res.val.total, res.err
}

func (r *pgRepo) GetDevices(ctx context.Context, ownerID string) ([]DeviceSummary, error) {
	resp := make(chan response[[]DeviceSummary], 1)
	r.getDevicesCh <- getDevicesReq{ctx: ctx, ownerID: ownerID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetDeviceDetail(ctx context.Context, deviceID string) (*DeviceDetail, error) {
	resp := make(chan response[*DeviceDetail], 1)
	r.getDeviceDetailCh <- getDeviceDetailReq{ctx: ctx, deviceID: deviceID, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetNotifications(ctx context.Context, input GetNotificationsInput) ([]Notification, int, int, error) {
	resp := make(chan response[notificationsResult], 1)
	r.getNotificationsCh <- getNotificationsReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val.notifications, res.val.unreadCount, res.val.total, res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doCreateProfile(ctx context.Context, input CreateProfileInput) (*Profile, error) {
	now := time.Now()
	var p Profile
	err := r.pool.QueryRow(ctx,
		`INSERT INTO scitizen_profiles (user_id, tos_accepted, tos_version, tos_accepted_at)
		 VALUES ($1, true, $2, $3)
		 RETURNING user_id, tos_accepted, tos_version, tos_accepted_at,
		           device_registered, campaign_enrolled, first_reading, created_at, updated_at`,
		input.UserID, input.TOSVersion, now,
	).Scan(&p.UserID, &p.TOSAccepted, &p.TOSVersion, &p.TOSAcceptedAt,
		&p.DeviceRegistered, &p.CampaignEnrolled, &p.FirstReading, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert scitizen profile: %w", err)
	}
	return &p, nil
}

func (r *pgRepo) doGetProfile(ctx context.Context, userID string) (*Profile, error) {
	var p Profile
	err := r.pool.QueryRow(ctx,
		`SELECT user_id, tos_accepted, tos_version, tos_accepted_at,
		        device_registered, campaign_enrolled, first_reading, created_at, updated_at
		 FROM scitizen_profiles WHERE user_id = $1`,
		userID,
	).Scan(&p.UserID, &p.TOSAccepted, &p.TOSVersion, &p.TOSAcceptedAt,
		&p.DeviceRegistered, &p.CampaignEnrolled, &p.FirstReading, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get scitizen profile: %w", err)
	}
	return &p, nil
}

func (r *pgRepo) doUpdateOnboarding(ctx context.Context, input UpdateOnboardingInput) error {
	query := `UPDATE scitizen_profiles SET updated_at = now()`
	args := []any{}
	argIdx := 1

	if input.DeviceRegistered != nil {
		query += fmt.Sprintf(", device_registered = $%d", argIdx)
		args = append(args, *input.DeviceRegistered)
		argIdx++
	}
	if input.CampaignEnrolled != nil {
		query += fmt.Sprintf(", campaign_enrolled = $%d", argIdx)
		args = append(args, *input.CampaignEnrolled)
		argIdx++
	}
	if input.FirstReading != nil {
		query += fmt.Sprintf(", first_reading = $%d", argIdx)
		args = append(args, *input.FirstReading)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE user_id = $%d", argIdx)
	args = append(args, input.UserID)

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update onboarding: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("scitizen profile %s not found", input.UserID)
	}
	return nil
}

func (r *pgRepo) doGetDashboard(ctx context.Context, userID string) (*Dashboard, error) {
	d := &Dashboard{}

	// Active enrollments count
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM campaign_enrollments WHERE scitizen_id = $1 AND status = 'active'`,
		userID,
	).Scan(&d.ActiveEnrollments)
	if err != nil {
		return nil, fmt.Errorf("count enrollments: %w", err)
	}

	// Reading counts from readings table
	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(COUNT(*), 0),
		        COALESCE(SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END), 0)
		 FROM readings r
		 JOIN devices d ON r.device_id = d.id
		 WHERE d.owner_id = $1`,
		userID,
	).Scan(&d.TotalReadings, &d.AcceptedReadings)
	if err != nil {
		return nil, fmt.Errorf("count readings: %w", err)
	}

	// Contribution score
	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(total, 0) FROM scitizen_scores WHERE scitizen_id = $1`,
		userID,
	).Scan(&d.ContributionScore)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("get score: %w", err)
	}

	// Badges
	badgeRows, err := r.pool.Query(ctx,
		`SELECT id, badge_type, awarded_at FROM badges WHERE scitizen_id = $1 ORDER BY awarded_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get badges: %w", err)
	}
	defer badgeRows.Close()
	for badgeRows.Next() {
		var b Badge
		if err := badgeRows.Scan(&b.ID, &b.BadgeType, &b.AwardedAt); err != nil {
			return nil, fmt.Errorf("scan badge: %w", err)
		}
		d.Badges = append(d.Badges, b)
	}
	if err := badgeRows.Err(); err != nil {
		return nil, err
	}

	// Enrollments
	enrollRows, err := r.pool.Query(ctx,
		`SELECT id, device_id, campaign_id, status, enrolled_at
		 FROM campaign_enrollments WHERE scitizen_id = $1 ORDER BY enrolled_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get enrollments: %w", err)
	}
	defer enrollRows.Close()
	for enrollRows.Next() {
		var e Enrollment
		if err := enrollRows.Scan(&e.ID, &e.DeviceID, &e.CampaignID, &e.Status, &e.EnrolledAt); err != nil {
			return nil, fmt.Errorf("scan enrollment: %w", err)
		}
		d.Enrollments = append(d.Enrollments, e)
	}
	return d, enrollRows.Err()
}

func (r *pgRepo) doGetContributions(ctx context.Context, userID string) ([]ReadingHistory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.device_id, r.campaign_id,
		        COUNT(*) as total,
		        SUM(CASE WHEN r.status = 'accepted' THEN 1 ELSE 0 END) as accepted,
		        SUM(CASE WHEN r.status = 'rejected' THEN 1 ELSE 0 END) as rejected
		 FROM readings r
		 JOIN devices d ON r.device_id = d.id
		 WHERE d.owner_id = $1
		 GROUP BY r.device_id, r.campaign_id
		 ORDER BY r.device_id, r.campaign_id`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get contributions: %w", err)
	}
	defer rows.Close()

	var results []ReadingHistory
	for rows.Next() {
		var h ReadingHistory
		if err := rows.Scan(&h.DeviceID, &h.CampaignID, &h.Total, &h.Accepted, &h.Rejected); err != nil {
			return nil, fmt.Errorf("scan reading history: %w", err)
		}
		results = append(results, h)
	}
	return results, rows.Err()
}

func (r *pgRepo) doBrowseCampaigns(ctx context.Context, input BrowseInput) (browseResult, error) {
	query := `SELECT c.id, c.status, c.window_start, c.window_end, c.created_at,
	                 COALESCE(e.cnt, 0) as enrollment_count
	          FROM campaigns c
	          LEFT JOIN (SELECT campaign_id, COUNT(*) as cnt FROM campaign_enrollments WHERE status = 'active' GROUP BY campaign_id) e
	            ON c.id = e.campaign_id
	          WHERE c.status = 'published'`
	args := []any{}
	argIdx := 1

	if input.SensorType != nil {
		query += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM campaign_eligibility ce
			WHERE ce.campaign_id = c.id AND $%d = ANY(ce.required_sensors)
		)`, argIdx)
		args = append(args, *input.SensorType)
		argIdx++
	}

	// Count total before pagination
	countQuery := `SELECT COUNT(*) FROM (` + query + `) sub`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return browseResult{}, fmt.Errorf("count campaigns: %w", err)
	}

	query += " ORDER BY c.created_at DESC"
	if input.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, input.Limit)
		argIdx++
	}
	if input.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, input.Offset)
		argIdx++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return browseResult{}, fmt.Errorf("browse campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []CampaignSummary
	for rows.Next() {
		var cs CampaignSummary
		if err := rows.Scan(&cs.ID, &cs.Status, &cs.WindowStart, &cs.WindowEnd, &cs.CreatedAt, &cs.EnrollmentCount); err != nil {
			return browseResult{}, fmt.Errorf("scan campaign summary: %w", err)
		}
		campaigns = append(campaigns, cs)
	}
	if err := rows.Err(); err != nil {
		return browseResult{}, err
	}

	// Fetch required sensors for each campaign
	for i := range campaigns {
		sensorRows, err := r.pool.Query(ctx,
			`SELECT DISTINCT unnest(required_sensors) FROM campaign_eligibility WHERE campaign_id = $1`,
			campaigns[i].ID,
		)
		if err != nil {
			return browseResult{}, fmt.Errorf("get sensors: %w", err)
		}
		for sensorRows.Next() {
			var s string
			if err := sensorRows.Scan(&s); err != nil {
				sensorRows.Close()
				return browseResult{}, fmt.Errorf("scan sensor: %w", err)
			}
			campaigns[i].RequiredSensors = append(campaigns[i].RequiredSensors, s)
		}
		sensorRows.Close()
	}

	return browseResult{campaigns: campaigns, total: total}, nil
}

func (r *pgRepo) doGetCampaignDetail(ctx context.Context, campaignID string) (*CampaignDetail, error) {
	d := &CampaignDetail{CampaignID: campaignID}

	err := r.pool.QueryRow(ctx,
		`SELECT status, window_start, window_end FROM campaigns WHERE id = $1`,
		campaignID,
	).Scan(&d.Status, &d.WindowStart, &d.WindowEnd)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign %s not found", campaignID)
		}
		return nil, fmt.Errorf("get campaign: %w", err)
	}

	// Enrollment count
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM campaign_enrollments WHERE campaign_id = $1 AND status = 'active'`,
		campaignID,
	).Scan(&d.EnrollmentCount)

	// Parameters
	paramRows, err := r.pool.Query(ctx,
		`SELECT name, unit, min_range, max_range, precision FROM campaign_parameters WHERE campaign_id = $1`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get parameters: %w", err)
	}
	defer paramRows.Close()
	for paramRows.Next() {
		var p Parameter
		if err := paramRows.Scan(&p.Name, &p.Unit, &p.MinRange, &p.MaxRange, &p.Precision); err != nil {
			return nil, fmt.Errorf("scan parameter: %w", err)
		}
		d.Parameters = append(d.Parameters, p)
	}

	// Regions
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
		d.Regions = append(d.Regions, reg)
	}

	// Eligibility
	eligRows, err := r.pool.Query(ctx,
		`SELECT device_class, tier, required_sensors, firmware_min FROM campaign_eligibility WHERE campaign_id = $1`,
		campaignID,
	)
	if err != nil {
		return nil, fmt.Errorf("get eligibility: %w", err)
	}
	defer eligRows.Close()
	for eligRows.Next() {
		var e EligibilityCriteria
		if err := eligRows.Scan(&e.DeviceClass, &e.Tier, &e.RequiredSensors, &e.FirmwareMin); err != nil {
			return nil, fmt.Errorf("scan eligibility: %w", err)
		}
		d.Eligibility = append(d.Eligibility, e)
	}

	return d, nil
}

func (r *pgRepo) doSearchCampaigns(ctx context.Context, input SearchInput) (browseResult, error) {
	query := `SELECT c.id, c.status, c.window_start, c.window_end, c.created_at,
	                 COALESCE(e.cnt, 0) as enrollment_count
	          FROM campaigns c
	          LEFT JOIN (SELECT campaign_id, COUNT(*) as cnt FROM campaign_enrollments WHERE status = 'active' GROUP BY campaign_id) e
	            ON c.id = e.campaign_id
	          WHERE c.status = 'published'
	            AND (c.id ILIKE $1
	                 OR EXISTS (SELECT 1 FROM campaign_parameters cp WHERE cp.campaign_id = c.id AND cp.name ILIKE $1))`
	searchPattern := "%" + input.Query + "%"
	args := []any{searchPattern}
	argIdx := 2

	countQuery := `SELECT COUNT(*) FROM (` + query + `) sub`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return browseResult{}, fmt.Errorf("count search results: %w", err)
	}

	query += " ORDER BY c.created_at DESC"
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
		return browseResult{}, fmt.Errorf("search campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []CampaignSummary
	for rows.Next() {
		var cs CampaignSummary
		if err := rows.Scan(&cs.ID, &cs.Status, &cs.WindowStart, &cs.WindowEnd, &cs.CreatedAt, &cs.EnrollmentCount); err != nil {
			return browseResult{}, fmt.Errorf("scan campaign: %w", err)
		}
		campaigns = append(campaigns, cs)
	}
	return browseResult{campaigns: campaigns, total: total}, rows.Err()
}

func (r *pgRepo) doGetDevices(ctx context.Context, ownerID string) ([]DeviceSummary, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT d.id, d.status, d.class, d.firmware_version, d.tier, d.sensors,
		        COALESCE(e.cnt, 0) as active_enrollments, d.last_seen
		 FROM devices d
		 LEFT JOIN (SELECT device_id, COUNT(*) as cnt FROM campaign_enrollments WHERE status = 'active' GROUP BY device_id) e
		   ON d.id = e.device_id
		 WHERE d.owner_id = $1
		 ORDER BY d.created_at DESC`,
		ownerID,
	)
	if err != nil {
		return nil, fmt.Errorf("get devices: %w", err)
	}
	defer rows.Close()

	var devices []DeviceSummary
	for rows.Next() {
		var ds DeviceSummary
		if err := rows.Scan(&ds.ID, &ds.Status, &ds.Class, &ds.FirmwareVersion, &ds.Tier, &ds.Sensors,
			&ds.ActiveEnrollments, &ds.LastSeen); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		devices = append(devices, ds)
	}
	return devices, rows.Err()
}

func (r *pgRepo) doGetDeviceDetail(ctx context.Context, deviceID string) (*DeviceDetail, error) {
	d := &DeviceDetail{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, owner_id, status, class, firmware_version, tier, sensors, cert_serial, created_at
		 FROM devices WHERE id = $1`,
		deviceID,
	).Scan(&d.ID, &d.OwnerID, &d.Status, &d.Class, &d.FirmwareVersion, &d.Tier, &d.Sensors, &d.CertSerial, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("device %s not found", deviceID)
		}
		return nil, fmt.Errorf("get device: %w", err)
	}

	// Enrollments
	enrollRows, err := r.pool.Query(ctx,
		`SELECT id, device_id, campaign_id, status, enrolled_at
		 FROM campaign_enrollments WHERE device_id = $1 ORDER BY enrolled_at DESC`,
		deviceID,
	)
	if err != nil {
		return nil, fmt.Errorf("get device enrollments: %w", err)
	}
	defer enrollRows.Close()
	for enrollRows.Next() {
		var e Enrollment
		if err := enrollRows.Scan(&e.ID, &e.DeviceID, &e.CampaignID, &e.Status, &e.EnrolledAt); err != nil {
			return nil, fmt.Errorf("scan enrollment: %w", err)
		}
		d.Enrollments = append(d.Enrollments, e)
	}

	return d, nil
}

func (r *pgRepo) doGetNotifications(ctx context.Context, input GetNotificationsInput) (notificationsResult, error) {
	// Unread count
	var unreadCount int
	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`,
		input.UserID,
	).Scan(&unreadCount)

	// Total and filtered query
	query := `SELECT id, user_id, type, message, read, resource_link, created_at
	          FROM notifications WHERE user_id = $1`
	args := []any{input.UserID}
	argIdx := 2

	if input.TypeFilter != nil {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, *input.TypeFilter)
		argIdx++
	}

	// Total count
	countQuery := `SELECT COUNT(*) FROM (` + query + `) sub`
	var total int
	r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)

	query += " ORDER BY created_at DESC"
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
		return notificationsResult{}, fmt.Errorf("get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.Read, &n.ResourceLink, &n.CreatedAt); err != nil {
			return notificationsResult{}, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	return notificationsResult{notifications: notifications, unreadCount: unreadCount, total: total}, rows.Err()
}

// Ensure ulid import is used.
var _ = ulid.Make
