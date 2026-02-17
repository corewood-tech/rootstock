package authorization

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/open-policy-agent/opa/v1/rego"

	"rootstock/web-server/global/observability"
)

//go:embed policies/authz.rego
var basePolicy string

type evalReq struct {
	ctx   context.Context
	input AuthzInput
	resp  chan evalResp
}

type evalResp struct {
	result *AuthzResult
	err    error
}

type recompileReq struct {
	ctx  context.Context
	resp chan error
}

type stopReq struct {
	resp chan struct{}
}

type opaState struct {
	prepared    rego.PreparedEvalQuery
	initialized bool
	logger      observability.Logger
	evalCh      chan evalReq
	recompileCh chan recompileReq
	stopCh      chan stopReq
}

// NewOPARepository creates a new authorization repository backed by OPA.
// It starts a manager goroutine that owns all mutable state.
func NewOPARepository() Repository {
	s := &opaState{
		logger:      observability.GetLogger("authorization"),
		evalCh:      make(chan evalReq),
		recompileCh: make(chan recompileReq),
		stopCh:      make(chan stopReq),
	}
	go s.manage()
	return s
}

func (s *opaState) manage() {
	for {
		select {
		case req := <-s.evalCh:
			result, err := s.doEval(req.ctx, req.input)
			req.resp <- evalResp{result: result, err: err}

		case req := <-s.recompileCh:
			req.resp <- s.doRecompile(req.ctx)

		case req := <-s.stopCh:
			close(req.resp)
			return
		}
	}
}

// Evaluate checks authorization for the given input.
func (s *opaState) Evaluate(ctx context.Context, input AuthzInput) (*AuthzResult, error) {
	resp := make(chan evalResp, 1)
	s.evalCh <- evalReq{ctx: ctx, input: input, resp: resp}
	r := <-resp
	return r.result, r.err
}

// Recompile prepares the OPA policy for evaluation.
func (s *opaState) Recompile(ctx context.Context) error {
	resp := make(chan error, 1)
	s.recompileCh <- recompileReq{ctx: ctx, resp: resp}
	return <-resp
}

func (s *opaState) doEval(ctx context.Context, input AuthzInput) (*AuthzResult, error) {
	if !s.initialized {
		return nil, fmt.Errorf("authorization not initialized — call Recompile first")
	}

	inputMap, err := structToMap(input)
	if err != nil {
		return nil, fmt.Errorf("convert input: %w", err)
	}

	start := time.Now()
	results, err := s.prepared.Eval(ctx, rego.EvalInput(inputMap))
	elapsed := time.Since(start)

	if err != nil {
		s.logDecision(ctx, input, nil, elapsed, err)
		return nil, fmt.Errorf("evaluate policy: %w", err)
	}

	if len(results) == 0 {
		err := fmt.Errorf("no result from policy evaluation")
		s.logDecision(ctx, input, nil, elapsed, err)
		return nil, err
	}

	decision, err := extractDecision(results)
	if err != nil {
		s.logDecision(ctx, input, nil, elapsed, err)
		return nil, fmt.Errorf("extract decision: %w", err)
	}

	s.logDecision(ctx, input, decision, elapsed, nil)
	return decision, nil
}

func (s *opaState) doRecompile(ctx context.Context) error {
	s.logger.Info(ctx, "recompiling authorization policy", nil)

	prepared, err := rego.New(
		rego.Query("data.authz.decision"),
		rego.Module("authz.rego", basePolicy),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("prepare policy: %w", err)
	}

	s.prepared = prepared
	s.initialized = true

	s.logger.Info(ctx, "authorization policy compiled", nil)
	return nil
}

func (s *opaState) logDecision(ctx context.Context, input AuthzInput, result *AuthzResult, elapsed time.Duration, evalErr error) {
	attrs := map[string]interface{}{
		"method":          input.Method,
		"session_user_id": input.SessionUserID,
		"duration_ms":     elapsed.Milliseconds(),
	}

	if evalErr != nil {
		attrs["error"] = evalErr.Error()
		s.logger.Error(ctx, "authz decision error", attrs)
		return
	}

	attrs["allow"] = result.Allow
	attrs["reason"] = result.Reason

	if result.Allow {
		s.logger.Info(ctx, "authz decision", attrs)
	} else {
		s.logger.Warn(ctx, "authz decision denied", attrs)
	}
}

func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func extractDecision(results rego.ResultSet) (*AuthzResult, error) {
	if len(results) == 0 || len(results[0].Expressions) == 0 {
		return nil, fmt.Errorf("empty result set")
	}

	decisionMap, ok := results[0].Expressions[0].Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("decision is not an object: %T", results[0].Expressions[0].Value)
	}

	allow, _ := decisionMap["allow"].(bool)
	reason, _ := decisionMap["reason"].(string)

	return &AuthzResult{
		Allow:  allow,
		Reason: reason,
	}, nil
}
