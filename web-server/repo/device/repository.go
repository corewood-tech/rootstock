package device

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type response[T any] struct {
	val T
	err error
}

type createReq struct {
	ctx   context.Context
	input CreateDeviceInput
	resp  chan response[*Device]
}

type getReq struct {
	ctx  context.Context
	id   string
	resp chan response[*Device]
}

type getCapsReq struct {
	ctx  context.Context
	id   string
	resp chan response[*DeviceCapabilities]
}

type updateStatusReq struct {
	ctx    context.Context
	id     string
	status string
	resp   chan response[struct{}]
}

type queryByClassReq struct {
	ctx   context.Context
	input QueryByClassInput
	resp  chan response[[]Device]
}

type genCodeReq struct {
	ctx   context.Context
	input GenerateCodeInput
	resp  chan response[*EnrollmentCode]
}

type redeemCodeReq struct {
	ctx  context.Context
	code string
	resp chan response[*EnrollmentCode]
}

type enrollReq struct {
	ctx        context.Context
	deviceID   string
	campaignID string
	resp       chan response[struct{}]
}

type shutdownReq struct {
	resp chan struct{}
}

type pgRepo struct {
	pool            *pgxpool.Pool
	createCh        chan createReq
	getCh           chan getReq
	getCapsCh       chan getCapsReq
	updateStatusCh  chan updateStatusReq
	queryByClassCh  chan queryByClassReq
	genCodeCh       chan genCodeReq
	redeemCodeCh    chan redeemCodeReq
	enrollCh        chan enrollReq
	shutdownCh      chan shutdownReq
}

// NewRepository creates a device repository backed by Postgres.
func NewRepository(pool *pgxpool.Pool) Repository {
	r := &pgRepo{
		pool:           pool,
		createCh:       make(chan createReq),
		getCh:          make(chan getReq),
		getCapsCh:      make(chan getCapsReq),
		updateStatusCh: make(chan updateStatusReq),
		queryByClassCh: make(chan queryByClassReq),
		genCodeCh:      make(chan genCodeReq),
		redeemCodeCh:   make(chan redeemCodeReq),
		enrollCh:       make(chan enrollReq),
		shutdownCh:     make(chan shutdownReq),
	}
	go r.manage()
	return r
}

func (r *pgRepo) manage() {
	for {
		select {
		case req := <-r.createCh:
			val, err := r.doCreate(req.ctx, req.input)
			req.resp <- response[*Device]{val: val, err: err}
		case req := <-r.getCh:
			val, err := r.doGet(req.ctx, req.id)
			req.resp <- response[*Device]{val: val, err: err}
		case req := <-r.getCapsCh:
			val, err := r.doGetCapabilities(req.ctx, req.id)
			req.resp <- response[*DeviceCapabilities]{val: val, err: err}
		case req := <-r.updateStatusCh:
			err := r.doUpdateStatus(req.ctx, req.id, req.status)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.queryByClassCh:
			val, err := r.doQueryByClass(req.ctx, req.input)
			req.resp <- response[[]Device]{val: val, err: err}
		case req := <-r.genCodeCh:
			val, err := r.doGenerateCode(req.ctx, req.input)
			req.resp <- response[*EnrollmentCode]{val: val, err: err}
		case req := <-r.redeemCodeCh:
			val, err := r.doRedeemCode(req.ctx, req.code)
			req.resp <- response[*EnrollmentCode]{val: val, err: err}
		case req := <-r.enrollCh:
			err := r.doEnroll(req.ctx, req.deviceID, req.campaignID)
			req.resp <- response[struct{}]{err: err}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *pgRepo) Create(ctx context.Context, input CreateDeviceInput) (*Device, error) {
	resp := make(chan response[*Device], 1)
	r.createCh <- createReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) Get(ctx context.Context, id string) (*Device, error) {
	resp := make(chan response[*Device], 1)
	r.getCh <- getReq{ctx: ctx, id: id, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GetCapabilities(ctx context.Context, id string) (*DeviceCapabilities, error) {
	resp := make(chan response[*DeviceCapabilities], 1)
	r.getCapsCh <- getCapsReq{ctx: ctx, id: id, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	resp := make(chan response[struct{}], 1)
	r.updateStatusCh <- updateStatusReq{ctx: ctx, id: id, status: status, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) QueryByClass(ctx context.Context, input QueryByClassInput) ([]Device, error) {
	resp := make(chan response[[]Device], 1)
	r.queryByClassCh <- queryByClassReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) GenerateEnrollmentCode(ctx context.Context, input GenerateCodeInput) (*EnrollmentCode, error) {
	resp := make(chan response[*EnrollmentCode], 1)
	r.genCodeCh <- genCodeReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) RedeemEnrollmentCode(ctx context.Context, code string) (*EnrollmentCode, error) {
	resp := make(chan response[*EnrollmentCode], 1)
	r.redeemCodeCh <- redeemCodeReq{ctx: ctx, code: code, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *pgRepo) EnrollInCampaign(ctx context.Context, deviceID string, campaignID string) error {
	resp := make(chan response[struct{}], 1)
	r.enrollCh <- enrollReq{ctx: ctx, deviceID: deviceID, campaignID: campaignID, resp: resp}
	res := <-resp
	return res.err
}

func (r *pgRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// --- implementation ---

func (r *pgRepo) doCreate(ctx context.Context, input CreateDeviceInput) (*Device, error) {
	var d Device
	err := r.pool.QueryRow(ctx,
		`INSERT INTO devices (owner_id, class, firmware_version, tier, sensors)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, owner_id, status, class, firmware_version, tier, sensors, cert_serial, created_at`,
		input.OwnerID, input.Class, input.FirmwareVersion, input.Tier, input.Sensors,
	).Scan(&d.ID, &d.OwnerID, &d.Status, &d.Class, &d.FirmwareVersion, &d.Tier, &d.Sensors, &d.CertSerial, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert device: %w", err)
	}
	return &d, nil
}

func (r *pgRepo) doGet(ctx context.Context, id string) (*Device, error) {
	var d Device
	err := r.pool.QueryRow(ctx,
		`SELECT id, owner_id, status, class, firmware_version, tier, sensors, cert_serial, created_at
		 FROM devices WHERE id = $1`,
		id,
	).Scan(&d.ID, &d.OwnerID, &d.Status, &d.Class, &d.FirmwareVersion, &d.Tier, &d.Sensors, &d.CertSerial, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("device %s not found", id)
		}
		return nil, fmt.Errorf("get device: %w", err)
	}
	return &d, nil
}

func (r *pgRepo) doGetCapabilities(ctx context.Context, id string) (*DeviceCapabilities, error) {
	var c DeviceCapabilities
	err := r.pool.QueryRow(ctx,
		`SELECT class, tier, sensors, firmware_version FROM devices WHERE id = $1`,
		id,
	).Scan(&c.Class, &c.Tier, &c.Sensors, &c.FirmwareVersion)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("device %s not found", id)
		}
		return nil, fmt.Errorf("get capabilities: %w", err)
	}
	return &c, nil
}

func (r *pgRepo) doUpdateStatus(ctx context.Context, id string, status string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE devices SET status = $1 WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("device %s not found", id)
	}
	return nil
}

func (r *pgRepo) doQueryByClass(ctx context.Context, input QueryByClassInput) ([]Device, error) {
	query := `SELECT id, owner_id, status, class, firmware_version, tier, sensors, cert_serial, created_at
	          FROM devices WHERE class = $1`
	args := []any{input.Class}
	argIdx := 2

	if input.FirmwareMinGte != "" {
		query += fmt.Sprintf(" AND firmware_version >= $%d", argIdx)
		args = append(args, input.FirmwareMinGte)
		argIdx++
	}
	if input.FirmwareMaxLte != "" {
		query += fmt.Sprintf(" AND firmware_version <= $%d", argIdx)
		args = append(args, input.FirmwareMaxLte)
		argIdx++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query by class: %w", err)
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.OwnerID, &d.Status, &d.Class, &d.FirmwareVersion, &d.Tier, &d.Sensors, &d.CertSerial, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

func (r *pgRepo) doGenerateCode(ctx context.Context, input GenerateCodeInput) (*EnrollmentCode, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(input.TTL) * time.Second)
	var ec EnrollmentCode
	err := r.pool.QueryRow(ctx,
		`INSERT INTO enrollment_codes (code, device_id, expires_at)
		 VALUES ($1, $2, $3)
		 RETURNING code, device_id, expires_at, used`,
		input.Code, input.DeviceID, expiresAt,
	).Scan(&ec.Code, &ec.DeviceID, &ec.ExpiresAt, &ec.Used)
	if err != nil {
		return nil, fmt.Errorf("insert enrollment code: %w", err)
	}
	return &ec, nil
}

func (r *pgRepo) doRedeemCode(ctx context.Context, code string) (*EnrollmentCode, error) {
	var ec EnrollmentCode
	err := r.pool.QueryRow(ctx,
		`UPDATE enrollment_codes
		 SET used = true
		 WHERE code = $1 AND used = false AND expires_at > now()
		 RETURNING code, device_id, expires_at, used`,
		code,
	).Scan(&ec.Code, &ec.DeviceID, &ec.ExpiresAt, &ec.Used)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("code %s not found, expired, or already used", code)
		}
		return nil, fmt.Errorf("redeem code: %w", err)
	}
	return &ec, nil
}

func (r *pgRepo) doEnroll(ctx context.Context, deviceID string, campaignID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO device_campaigns (device_id, campaign_id) VALUES ($1, $2)`,
		deviceID, campaignID,
	)
	if err != nil {
		return fmt.Errorf("enroll device: %w", err)
	}
	return nil
}
