package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type response[T any] struct {
	val T
	err error
}

// --- Request types for channel-based concurrency ---

type getCurrentStateReq struct {
	ctx         context.Context
	campaignRef string
	resp        chan response[*CampaignInstanceState]
}

type transitionStateReq struct {
	ctx   context.Context
	input TransitionInput
	resp  chan response[*CampaignInstanceState]
}

type getValidTransitionsReq struct {
	ctx         context.Context
	campaignRef string
	resp        chan response[[]ValidTransition]
}

type initCampaignStateReq struct {
	ctx         context.Context
	campaignRef string
	resp        chan response[*CampaignInstanceState]
}

type getBaselineReq struct {
	ctx           context.Context
	campaignRef   string
	parameterName string
	resp          chan response[*AnomalyBaseline]
}

type updateBaselineReq struct {
	ctx   context.Context
	input UpdateBaselineInput
	resp  chan response[*AnomalyBaseline]
}

type checkAnomalyReq struct {
	ctx   context.Context
	input CheckAnomalyInput
	resp  chan response[*AnomalyFlag]
}

type addEnrollmentReq struct {
	ctx   context.Context
	input EnrollmentInput
	resp  chan response[struct{}]
}

type withdrawEnrollmentReq struct {
	ctx         context.Context
	deviceRef   string
	campaignRef string
	resp        chan response[struct{}]
}

type getDeviceCampaignsReq struct {
	ctx       context.Context
	deviceRef string
	resp      chan response[[]EnrollmentEdge]
}

type getCampaignDevicesReq struct {
	ctx         context.Context
	campaignRef string
	resp        chan response[[]EnrollmentEdge]
}

type getSharedDeviceCampaignsReq struct {
	ctx         context.Context
	campaignRef string
	resp        chan response[[]string]
}

type shutdownReq struct {
	resp chan struct{}
}

type dgraphRepo struct {
	client *dgo.Dgraph
	conn   *grpc.ClientConn

	getCurrentStateCh         chan getCurrentStateReq
	transitionStateCh         chan transitionStateReq
	getValidTransitionsCh     chan getValidTransitionsReq
	initCampaignStateCh       chan initCampaignStateReq
	getBaselineCh             chan getBaselineReq
	updateBaselineCh          chan updateBaselineReq
	checkAnomalyCh            chan checkAnomalyReq
	addEnrollmentCh           chan addEnrollmentReq
	withdrawEnrollmentCh      chan withdrawEnrollmentReq
	getDeviceCampaignsCh      chan getDeviceCampaignsReq
	getCampaignDevicesCh      chan getCampaignDevicesReq
	getSharedDeviceCampaignsCh chan getSharedDeviceCampaignsReq
	shutdownCh                chan shutdownReq
}

// NewDgraphRepository creates a graph repository backed by Dgraph.
func NewDgraphRepository(alphaAddr string) (Repository, error) {
	conn, err := grpc.NewClient(alphaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to dgraph alpha at %s: %w", alphaAddr, err)
	}

	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	r := &dgraphRepo{
		client:                     client,
		conn:                       conn,
		getCurrentStateCh:          make(chan getCurrentStateReq),
		transitionStateCh:          make(chan transitionStateReq),
		getValidTransitionsCh:      make(chan getValidTransitionsReq),
		initCampaignStateCh:        make(chan initCampaignStateReq),
		getBaselineCh:              make(chan getBaselineReq),
		updateBaselineCh:           make(chan updateBaselineReq),
		checkAnomalyCh:             make(chan checkAnomalyReq),
		addEnrollmentCh:            make(chan addEnrollmentReq),
		withdrawEnrollmentCh:       make(chan withdrawEnrollmentReq),
		getDeviceCampaignsCh:       make(chan getDeviceCampaignsReq),
		getCampaignDevicesCh:       make(chan getCampaignDevicesReq),
		getSharedDeviceCampaignsCh: make(chan getSharedDeviceCampaignsReq),
		shutdownCh:                 make(chan shutdownReq),
	}
	go r.manage()
	return r, nil
}

func (r *dgraphRepo) manage() {
	for {
		select {
		case req := <-r.getCurrentStateCh:
			val, err := r.doGetCurrentState(req.ctx, req.campaignRef)
			req.resp <- response[*CampaignInstanceState]{val, err}

		case req := <-r.transitionStateCh:
			val, err := r.doTransitionState(req.ctx, req.input)
			req.resp <- response[*CampaignInstanceState]{val, err}

		case req := <-r.getValidTransitionsCh:
			val, err := r.doGetValidTransitions(req.ctx, req.campaignRef)
			req.resp <- response[[]ValidTransition]{val, err}

		case req := <-r.initCampaignStateCh:
			val, err := r.doInitCampaignState(req.ctx, req.campaignRef)
			req.resp <- response[*CampaignInstanceState]{val, err}

		case req := <-r.getBaselineCh:
			val, err := r.doGetBaseline(req.ctx, req.campaignRef, req.parameterName)
			req.resp <- response[*AnomalyBaseline]{val, err}

		case req := <-r.updateBaselineCh:
			val, err := r.doUpdateBaseline(req.ctx, req.input)
			req.resp <- response[*AnomalyBaseline]{val, err}

		case req := <-r.checkAnomalyCh:
			val, err := r.doCheckAnomaly(req.ctx, req.input)
			req.resp <- response[*AnomalyFlag]{val, err}

		case req := <-r.addEnrollmentCh:
			err := r.doAddEnrollment(req.ctx, req.input)
			req.resp <- response[struct{}]{struct{}{}, err}

		case req := <-r.withdrawEnrollmentCh:
			err := r.doWithdrawEnrollment(req.ctx, req.deviceRef, req.campaignRef)
			req.resp <- response[struct{}]{struct{}{}, err}

		case req := <-r.getDeviceCampaignsCh:
			val, err := r.doGetDeviceCampaigns(req.ctx, req.deviceRef)
			req.resp <- response[[]EnrollmentEdge]{val, err}

		case req := <-r.getCampaignDevicesCh:
			val, err := r.doGetCampaignDevices(req.ctx, req.campaignRef)
			req.resp <- response[[]EnrollmentEdge]{val, err}

		case req := <-r.getSharedDeviceCampaignsCh:
			val, err := r.doGetSharedDeviceCampaigns(req.ctx, req.campaignRef)
			req.resp <- response[[]string]{val, err}

		case req := <-r.shutdownCh:
			r.conn.Close()
			req.resp <- struct{}{}
			return
		}
	}
}

// --- Public methods (send to channel, wait for response) ---

func (r *dgraphRepo) GetCurrentState(ctx context.Context, campaignRef string) (*CampaignInstanceState, error) {
	resp := make(chan response[*CampaignInstanceState], 1)
	r.getCurrentStateCh <- getCurrentStateReq{ctx, campaignRef, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) TransitionState(ctx context.Context, input TransitionInput) (*CampaignInstanceState, error) {
	resp := make(chan response[*CampaignInstanceState], 1)
	r.transitionStateCh <- transitionStateReq{ctx, input, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) GetValidTransitions(ctx context.Context, campaignRef string) ([]ValidTransition, error) {
	resp := make(chan response[[]ValidTransition], 1)
	r.getValidTransitionsCh <- getValidTransitionsReq{ctx, campaignRef, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) InitCampaignState(ctx context.Context, campaignRef string) (*CampaignInstanceState, error) {
	resp := make(chan response[*CampaignInstanceState], 1)
	r.initCampaignStateCh <- initCampaignStateReq{ctx, campaignRef, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) GetBaseline(ctx context.Context, campaignRef string, parameterName string) (*AnomalyBaseline, error) {
	resp := make(chan response[*AnomalyBaseline], 1)
	r.getBaselineCh <- getBaselineReq{ctx, campaignRef, parameterName, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) UpdateBaseline(ctx context.Context, input UpdateBaselineInput) (*AnomalyBaseline, error) {
	resp := make(chan response[*AnomalyBaseline], 1)
	r.updateBaselineCh <- updateBaselineReq{ctx, input, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) CheckAnomaly(ctx context.Context, input CheckAnomalyInput) (*AnomalyFlag, error) {
	resp := make(chan response[*AnomalyFlag], 1)
	r.checkAnomalyCh <- checkAnomalyReq{ctx, input, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) AddEnrollment(ctx context.Context, input EnrollmentInput) error {
	resp := make(chan response[struct{}], 1)
	r.addEnrollmentCh <- addEnrollmentReq{ctx, input, resp}
	res := <-resp
	return res.err
}

func (r *dgraphRepo) WithdrawEnrollment(ctx context.Context, deviceRef string, campaignRef string) error {
	resp := make(chan response[struct{}], 1)
	r.withdrawEnrollmentCh <- withdrawEnrollmentReq{ctx, deviceRef, campaignRef, resp}
	res := <-resp
	return res.err
}

func (r *dgraphRepo) GetDeviceCampaigns(ctx context.Context, deviceRef string) ([]EnrollmentEdge, error) {
	resp := make(chan response[[]EnrollmentEdge], 1)
	r.getDeviceCampaignsCh <- getDeviceCampaignsReq{ctx, deviceRef, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) GetCampaignDevices(ctx context.Context, campaignRef string) ([]EnrollmentEdge, error) {
	resp := make(chan response[[]EnrollmentEdge], 1)
	r.getCampaignDevicesCh <- getCampaignDevicesReq{ctx, campaignRef, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) GetSharedDeviceCampaigns(ctx context.Context, campaignRef string) ([]string, error) {
	resp := make(chan response[[]string], 1)
	r.getSharedDeviceCampaignsCh <- getSharedDeviceCampaignsReq{ctx, campaignRef, resp}
	res := <-resp
	return res.val, res.err
}

func (r *dgraphRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp}
	<-resp
}

// ============================================================
// Internal implementations
// ============================================================

// --- Campaign State Machine ---

type instanceQuery struct {
	Instances []struct {
		UID          string `json:"uid"`
		CampaignRef  string `json:"campaign_ref"`
		CurrentState []struct {
			UID       string `json:"uid"`
			StateName string `json:"state_name"`
		} `json:"current_state"`
		StateEnteredAt string `json:"state_entered_at"`
	} `json:"instances"`
}

func (r *dgraphRepo) doGetCurrentState(ctx context.Context, campaignRef string) (*CampaignInstanceState, error) {
	q := `query GetState($ref: string) {
		instances(func: eq(campaign_ref, $ref)) @filter(type(CampaignInstance)) {
			uid
			campaign_ref
			current_state {
				uid
				state_name
			}
			state_entered_at
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$ref": campaignRef})
	if err != nil {
		return nil, fmt.Errorf("query campaign state: %w", err)
	}

	var result instanceQuery
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal campaign state: %w", err)
	}

	if len(result.Instances) == 0 {
		return nil, fmt.Errorf("campaign instance not found: %s", campaignRef)
	}

	inst := result.Instances[0]
	stateName := ""
	if len(inst.CurrentState) > 0 {
		stateName = inst.CurrentState[0].StateName
	}

	enteredAt, _ := time.Parse(time.RFC3339, inst.StateEnteredAt)

	return &CampaignInstanceState{
		CampaignRef: campaignRef,
		StateName:   stateName,
		EnteredAt:   enteredAt,
	}, nil
}

func (r *dgraphRepo) doGetValidTransitions(ctx context.Context, campaignRef string) ([]ValidTransition, error) {
	// First get current state
	state, err := r.doGetCurrentState(ctx, campaignRef)
	if err != nil {
		return nil, err
	}

	q := `query Transitions($state: string) {
		transitions(func: type(Transition)) @filter(eq(val(from_state), $state)) @cascade {
			transition_from @filter(eq(state_name, $state)) {
				state_name
			}
			transition_to {
				state_name
			}
			transition_event {
				event_name
			}
			guard
			side_effect
		}
	}`

	// Simpler approach: query transitions whose from-state matches
	q = `query Transitions($state: string) {
		var(func: eq(state_name, $state)) @filter(type(CampaignState)) {
			from as uid
		}
		transitions(func: type(Transition)) @filter(uid_in(transition_from, uid(from))) {
			transition_to {
				state_name
			}
			transition_event {
				event_name
			}
			guard
			side_effect
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$state": state.StateName})
	if err != nil {
		return nil, fmt.Errorf("query transitions: %w", err)
	}

	var result struct {
		Transitions []struct {
			TransitionTo []struct {
				StateName string `json:"state_name"`
			} `json:"transition_to"`
			TransitionEvent []struct {
				EventName string `json:"event_name"`
			} `json:"transition_event"`
			Guard      string `json:"guard"`
			SideEffect string `json:"side_effect"`
		} `json:"transitions"`
	}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal transitions: %w", err)
	}

	var out []ValidTransition
	for _, t := range result.Transitions {
		vt := ValidTransition{
			Guard:      t.Guard,
			SideEffect: t.SideEffect,
		}
		if len(t.TransitionTo) > 0 {
			vt.TargetState = t.TransitionTo[0].StateName
		}
		if len(t.TransitionEvent) > 0 {
			vt.EventName = t.TransitionEvent[0].EventName
		}
		out = append(out, vt)
	}

	return out, nil
}

func (r *dgraphRepo) doTransitionState(ctx context.Context, input TransitionInput) (*CampaignInstanceState, error) {
	// Get valid transitions and find the one matching the event
	transitions, err := r.doGetValidTransitions(ctx, input.CampaignRef)
	if err != nil {
		return nil, err
	}

	var target *ValidTransition
	for i, t := range transitions {
		if t.EventName == input.EventName {
			target = &transitions[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("no valid transition for event %q from current state", input.EventName)
	}

	// Find the target state UID
	q := `query FindState($name: string) {
		state(func: eq(state_name, $name)) @filter(type(CampaignState)) {
			uid
		}
	}`
	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$name": target.TargetState})
	if err != nil {
		return nil, fmt.Errorf("find target state: %w", err)
	}

	var stateResult struct {
		State []struct {
			UID string `json:"uid"`
		} `json:"state"`
	}
	if err := json.Unmarshal(resp.Json, &stateResult); err != nil {
		return nil, fmt.Errorf("unmarshal target state: %w", err)
	}
	if len(stateResult.State) == 0 {
		return nil, fmt.Errorf("target state node not found: %s", target.TargetState)
	}
	targetUID := stateResult.State[0].UID

	// Find the campaign instance UID
	instQ := `query FindInstance($ref: string) {
		inst(func: eq(campaign_ref, $ref)) @filter(type(CampaignInstance)) {
			uid
		}
	}`
	resp2, err := txn.QueryWithVars(ctx, instQ, map[string]string{"$ref": input.CampaignRef})
	if err != nil {
		return nil, fmt.Errorf("find campaign instance: %w", err)
	}

	var instResult struct {
		Inst []struct {
			UID string `json:"uid"`
		} `json:"inst"`
	}
	if err := json.Unmarshal(resp2.Json, &instResult); err != nil {
		return nil, fmt.Errorf("unmarshal instance: %w", err)
	}
	if len(instResult.Inst) == 0 {
		return nil, fmt.Errorf("campaign instance not found: %s", input.CampaignRef)
	}
	instUID := instResult.Inst[0].UID

	// Mutate: update current_state and state_entered_at
	now := time.Now().UTC()
	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(
			`<%s> <current_state> <%s> .
			 <%s> <state_entered_at> "%s"^^<xs:dateTime> .`,
			instUID, targetUID,
			instUID, now.Format(time.RFC3339),
		)),
	}

	writeTxn := r.client.NewTxn()
	defer writeTxn.Discard(ctx)
	if _, err := writeTxn.Mutate(ctx, mu); err != nil {
		return nil, fmt.Errorf("transition state: %w", err)
	}
	if err := writeTxn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transition: %w", err)
	}

	return &CampaignInstanceState{
		CampaignRef: input.CampaignRef,
		StateName:   target.TargetState,
		EnteredAt:   now,
	}, nil
}

func (r *dgraphRepo) doInitCampaignState(ctx context.Context, campaignRef string) (*CampaignInstanceState, error) {
	// Find the "draft" state node
	q := `query FindDraft {
		state(func: eq(state_name, "draft")) @filter(type(CampaignState)) {
			uid
		}
	}`
	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("find draft state: %w", err)
	}

	var result struct {
		State []struct {
			UID string `json:"uid"`
		} `json:"state"`
	}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal draft state: %w", err)
	}
	if len(result.State) == 0 {
		return nil, fmt.Errorf("draft state node not found — seed the state machine first")
	}
	draftUID := result.State[0].UID

	now := time.Now().UTC()
	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(
			`_:inst <dgraph.type> "CampaignInstance" .
			 _:inst <campaign_ref> "%s" .
			 _:inst <current_state> <%s> .
			 _:inst <state_entered_at> "%s"^^<xs:dateTime> .`,
			campaignRef, draftUID, now.Format(time.RFC3339),
		)),
	}

	writeTxn := r.client.NewTxn()
	defer writeTxn.Discard(ctx)
	if _, err := writeTxn.Mutate(ctx, mu); err != nil {
		return nil, fmt.Errorf("init campaign state: %w", err)
	}
	if err := writeTxn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit init: %w", err)
	}

	return &CampaignInstanceState{
		CampaignRef: campaignRef,
		StateName:   "draft",
		EnteredAt:   now,
	}, nil
}

// --- Anomaly Detection ---

type baselineQuery struct {
	Baselines []struct {
		UID                 string  `json:"uid"`
		CampaignRef         string  `json:"campaign_ref"`
		ParameterName       string  `json:"parameter_name"`
		SampleCount         int64   `json:"sample_count"`
		RollingMean         float64 `json:"rolling_mean"`
		RollingM2           float64 `json:"rolling_m2"`
		RollingMin          float64 `json:"rolling_min"`
		RollingMax          float64 `json:"rolling_max"`
		BoundStddevMult     float64 `json:"bound_stddev_multiplier"`
		BoundHardMin        float64 `json:"bound_hard_min"`
		BoundHardMax        float64 `json:"bound_hard_max"`
		LastUpdated         string  `json:"last_updated"`
	} `json:"baselines"`
}

func (r *dgraphRepo) doGetBaseline(ctx context.Context, campaignRef, parameterName string) (*AnomalyBaseline, error) {
	q := `query GetBaseline($campaign: string, $param: string) {
		baselines(func: eq(campaign_ref, $campaign)) @filter(type(AnomalyBaseline) AND eq(parameter_name, $param)) {
			uid
			campaign_ref
			parameter_name
			sample_count
			rolling_mean
			rolling_m2
			rolling_min
			rolling_max
			bound_stddev_multiplier
			bound_hard_min
			bound_hard_max
			last_updated
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$campaign": campaignRef,
		"$param":    parameterName,
	})
	if err != nil {
		return nil, fmt.Errorf("query baseline: %w", err)
	}

	var result baselineQuery
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal baseline: %w", err)
	}

	if len(result.Baselines) == 0 {
		return nil, nil // no baseline yet
	}

	b := result.Baselines[0]
	lastUpdated, _ := time.Parse(time.RFC3339, b.LastUpdated)

	return &AnomalyBaseline{
		CampaignRef:      b.CampaignRef,
		ParameterName:    b.ParameterName,
		SampleCount:      b.SampleCount,
		Mean:             b.RollingMean,
		M2:               b.RollingM2,
		Min:              b.RollingMin,
		Max:              b.RollingMax,
		StddevMultiplier: b.BoundStddevMult,
		HardMin:          b.BoundHardMin,
		HardMax:          b.BoundHardMax,
		LastUpdated:      lastUpdated,
	}, nil
}

func (r *dgraphRepo) doUpdateBaseline(ctx context.Context, input UpdateBaselineInput) (*AnomalyBaseline, error) {
	existing, err := r.doGetBaseline(ctx, input.CampaignRef, input.ParameterName)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	if existing == nil {
		// First reading — create baseline node
		mu := &api.Mutation{
			SetNquads: []byte(fmt.Sprintf(
				`_:bl <dgraph.type> "AnomalyBaseline" .
				 _:bl <campaign_ref> "%s" .
				 _:bl <parameter_name> "%s" .
				 _:bl <sample_count> "1"^^<xs:int> .
				 _:bl <rolling_mean> "%f"^^<xs:float> .
				 _:bl <rolling_m2> "0"^^<xs:float> .
				 _:bl <rolling_min> "%f"^^<xs:float> .
				 _:bl <rolling_max> "%f"^^<xs:float> .
				 _:bl <bound_stddev_multiplier> "3"^^<xs:float> .
				 _:bl <last_updated> "%s"^^<xs:dateTime> .`,
				input.CampaignRef, input.ParameterName,
				input.Value, input.Value, input.Value,
				now.Format(time.RFC3339),
			)),
		}

		txn := r.client.NewTxn()
		defer txn.Discard(ctx)
		if _, err := txn.Mutate(ctx, mu); err != nil {
			return nil, fmt.Errorf("create baseline: %w", err)
		}
		if err := txn.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit baseline: %w", err)
		}

		return &AnomalyBaseline{
			CampaignRef:      input.CampaignRef,
			ParameterName:    input.ParameterName,
			SampleCount:      1,
			Mean:             input.Value,
			M2:               0,
			Min:              input.Value,
			Max:              input.Value,
			StddevMultiplier: 3.0,
			LastUpdated:      now,
		}, nil
	}

	// Welford's online algorithm
	n := existing.SampleCount + 1
	delta := input.Value - existing.Mean
	newMean := existing.Mean + delta/float64(n)
	delta2 := input.Value - newMean
	newM2 := existing.M2 + delta*delta2
	newMin := math.Min(existing.Min, input.Value)
	newMax := math.Max(existing.Max, input.Value)

	// Find baseline UID
	q := `query FindBL($campaign: string, $param: string) {
		bl(func: eq(campaign_ref, $campaign)) @filter(type(AnomalyBaseline) AND eq(parameter_name, $param)) {
			uid
		}
	}`
	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$campaign": input.CampaignRef,
		"$param":    input.ParameterName,
	})
	if err != nil {
		return nil, fmt.Errorf("find baseline uid: %w", err)
	}

	var blResult struct {
		BL []struct {
			UID string `json:"uid"`
		} `json:"bl"`
	}
	if err := json.Unmarshal(resp.Json, &blResult); err != nil {
		return nil, fmt.Errorf("unmarshal baseline uid: %w", err)
	}
	if len(blResult.BL) == 0 {
		return nil, fmt.Errorf("baseline uid not found")
	}
	blUID := blResult.BL[0].UID

	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(
			`<%s> <sample_count> "%d"^^<xs:int> .
			 <%s> <rolling_mean> "%f"^^<xs:float> .
			 <%s> <rolling_m2> "%f"^^<xs:float> .
			 <%s> <rolling_min> "%f"^^<xs:float> .
			 <%s> <rolling_max> "%f"^^<xs:float> .
			 <%s> <last_updated> "%s"^^<xs:dateTime> .`,
			blUID, n,
			blUID, newMean,
			blUID, newM2,
			blUID, newMin,
			blUID, newMax,
			blUID, now.Format(time.RFC3339),
		)),
	}

	writeTxn := r.client.NewTxn()
	defer writeTxn.Discard(ctx)
	if _, err := writeTxn.Mutate(ctx, mu); err != nil {
		return nil, fmt.Errorf("update baseline: %w", err)
	}
	if err := writeTxn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit baseline update: %w", err)
	}

	return &AnomalyBaseline{
		CampaignRef:      input.CampaignRef,
		ParameterName:    input.ParameterName,
		SampleCount:      n,
		Mean:             newMean,
		M2:               newM2,
		Min:              newMin,
		Max:              newMax,
		StddevMultiplier: existing.StddevMultiplier,
		HardMin:          existing.HardMin,
		HardMax:          existing.HardMax,
		LastUpdated:      now,
	}, nil
}

func (r *dgraphRepo) doCheckAnomaly(ctx context.Context, input CheckAnomalyInput) (*AnomalyFlag, error) {
	baseline, err := r.doGetBaseline(ctx, input.CampaignRef, input.ParameterName)
	if err != nil {
		return nil, err
	}
	if baseline == nil || baseline.SampleCount < 30 {
		// Not enough data to flag anomalies
		return nil, nil
	}

	variance := baseline.M2 / float64(baseline.SampleCount)
	stddev := math.Sqrt(variance)
	multiplier := baseline.StddevMultiplier
	if multiplier == 0 {
		multiplier = 3.0
	}

	lowerBound := baseline.Mean - multiplier*stddev
	upperBound := baseline.Mean + multiplier*stddev

	// Also check hard bounds if set
	if baseline.HardMin != 0 || baseline.HardMax != 0 {
		if baseline.HardMin != 0 && lowerBound < baseline.HardMin {
			lowerBound = baseline.HardMin
		}
		if baseline.HardMax != 0 && upperBound > baseline.HardMax {
			upperBound = baseline.HardMax
		}
	}

	if input.Value < lowerBound || input.Value > upperBound {
		return &AnomalyFlag{
			Reason:     fmt.Sprintf("value %.4f outside bounds [%.4f, %.4f] (mean=%.4f, stddev=%.4f, multiplier=%.1f)", input.Value, lowerBound, upperBound, baseline.Mean, stddev, multiplier),
			Value:      input.Value,
			LowerBound: lowerBound,
			UpperBound: upperBound,
			Mean:       baseline.Mean,
			Stddev:     stddev,
		}, nil
	}

	return nil, nil
}

// --- Device-Campaign Relationships ---

func (r *dgraphRepo) doAddEnrollment(ctx context.Context, input EnrollmentInput) error {
	// Upsert device node
	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(
			`_:enroll <dgraph.type> "Enrollment" .
			 _:enroll <device_ref> "%s" .
			 _:enroll <campaign_ref> "%s" .
			 _:enroll <owner_ref> "%s" .
			 _:enroll <enrolled_at> "%s"^^<xs:dateTime> .
			 _:enroll <enrollment_status> "active" .`,
			input.DeviceRef, input.CampaignRef, input.OwnerRef,
			input.EnrolledAt.Format(time.RFC3339),
		)),
	}

	txn := r.client.NewTxn()
	defer txn.Discard(ctx)
	if _, err := txn.Mutate(ctx, mu); err != nil {
		return fmt.Errorf("add enrollment: %w", err)
	}
	return txn.Commit(ctx)
}

func (r *dgraphRepo) doWithdrawEnrollment(ctx context.Context, deviceRef, campaignRef string) error {
	// Find the enrollment edge
	q := `query FindEnrollment($device: string, $campaign: string) {
		enrollment(func: eq(device_ref, $device)) @filter(type(Enrollment) AND eq(campaign_ref, $campaign) AND eq(enrollment_status, "active")) {
			uid
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$device":   deviceRef,
		"$campaign": campaignRef,
	})
	if err != nil {
		return fmt.Errorf("find enrollment: %w", err)
	}

	var result struct {
		Enrollment []struct {
			UID string `json:"uid"`
		} `json:"enrollment"`
	}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return fmt.Errorf("unmarshal enrollment: %w", err)
	}
	if len(result.Enrollment) == 0 {
		return fmt.Errorf("active enrollment not found for device %s in campaign %s", deviceRef, campaignRef)
	}

	now := time.Now().UTC()
	uid := result.Enrollment[0].UID

	mu := &api.Mutation{
		SetNquads: []byte(fmt.Sprintf(
			`<%s> <enrollment_status> "withdrawn" .
			 <%s> <withdrawn_at> "%s"^^<xs:dateTime> .`,
			uid, uid, now.Format(time.RFC3339),
		)),
	}

	writeTxn := r.client.NewTxn()
	defer writeTxn.Discard(ctx)
	if _, err := writeTxn.Mutate(ctx, mu); err != nil {
		return fmt.Errorf("withdraw enrollment: %w", err)
	}
	return writeTxn.Commit(ctx)
}

func (r *dgraphRepo) doGetDeviceCampaigns(ctx context.Context, deviceRef string) ([]EnrollmentEdge, error) {
	q := `query DeviceCampaigns($device: string) {
		enrollments(func: eq(device_ref, $device)) @filter(type(Enrollment) AND eq(enrollment_status, "active")) {
			device_ref
			campaign_ref
			owner_ref
			enrolled_at
			enrollment_status
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$device": deviceRef})
	if err != nil {
		return nil, fmt.Errorf("query device campaigns: %w", err)
	}

	var result struct {
		Enrollments []struct {
			DeviceRef        string `json:"device_ref"`
			CampaignRef      string `json:"campaign_ref"`
			OwnerRef         string `json:"owner_ref"`
			EnrolledAt       string `json:"enrolled_at"`
			EnrollmentStatus string `json:"enrollment_status"`
		} `json:"enrollments"`
	}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal enrollments: %w", err)
	}

	var out []EnrollmentEdge
	for _, e := range result.Enrollments {
		enrolledAt, _ := time.Parse(time.RFC3339, e.EnrolledAt)
		out = append(out, EnrollmentEdge{
			DeviceRef:        e.DeviceRef,
			CampaignRef:      e.CampaignRef,
			OwnerRef:         e.OwnerRef,
			EnrolledAt:       enrolledAt,
			EnrollmentStatus: e.EnrollmentStatus,
		})
	}

	return out, nil
}

func (r *dgraphRepo) doGetCampaignDevices(ctx context.Context, campaignRef string) ([]EnrollmentEdge, error) {
	q := `query CampaignDevices($campaign: string) {
		enrollments(func: eq(campaign_ref, $campaign)) @filter(type(Enrollment) AND eq(enrollment_status, "active")) {
			device_ref
			campaign_ref
			owner_ref
			enrolled_at
			enrollment_status
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$campaign": campaignRef})
	if err != nil {
		return nil, fmt.Errorf("query campaign devices: %w", err)
	}

	var result struct {
		Enrollments []struct {
			DeviceRef        string `json:"device_ref"`
			CampaignRef      string `json:"campaign_ref"`
			OwnerRef         string `json:"owner_ref"`
			EnrolledAt       string `json:"enrolled_at"`
			EnrollmentStatus string `json:"enrollment_status"`
		} `json:"enrollments"`
	}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal enrollments: %w", err)
	}

	var out []EnrollmentEdge
	for _, e := range result.Enrollments {
		enrolledAt, _ := time.Parse(time.RFC3339, e.EnrolledAt)
		out = append(out, EnrollmentEdge{
			DeviceRef:        e.DeviceRef,
			CampaignRef:      e.CampaignRef,
			OwnerRef:         e.OwnerRef,
			EnrolledAt:       enrolledAt,
			EnrollmentStatus: e.EnrollmentStatus,
		})
	}

	return out, nil
}

func (r *dgraphRepo) doGetSharedDeviceCampaigns(ctx context.Context, campaignRef string) ([]string, error) {
	// Two-hop traversal: campaign → enrolled devices → their other campaigns
	q := `query SharedCampaigns($campaign: string) {
		var(func: eq(campaign_ref, $campaign)) @filter(type(Enrollment) AND eq(enrollment_status, "active")) {
			devices as device_ref
		}
		shared(func: eq(device_ref, val(devices))) @filter(type(Enrollment) AND eq(enrollment_status, "active") AND NOT eq(campaign_ref, $campaign)) {
			campaign_ref
		}
	}`

	txn := r.client.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$campaign": campaignRef})
	if err != nil {
		return nil, fmt.Errorf("query shared campaigns: %w", err)
	}

	var result struct {
		Shared []struct {
			CampaignRef string `json:"campaign_ref"`
		} `json:"shared"`
	}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, fmt.Errorf("unmarshal shared: %w", err)
	}

	seen := make(map[string]bool)
	var out []string
	for _, s := range result.Shared {
		if !seen[s.CampaignRef] {
			seen[s.CampaignRef] = true
			out = append(out, s.CampaignRef)
		}
	}

	return out, nil
}
