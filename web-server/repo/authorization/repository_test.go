package authorization

import (
	"context"
	"testing"

	"github.com/open-policy-agent/opa/v1/rego"
)

// prepareTestPolicy compiles the embedded Rego policy for direct evaluation in tests,
// bypassing the goroutine manager so tests are self-contained.
func prepareTestPolicy(t *testing.T) rego.PreparedEvalQuery {
	t.Helper()
	prepared, err := rego.New(
		rego.Query("data.authz.decision"),
		rego.Module("authz.rego", basePolicy),
	).PrepareForEval(context.Background())
	if err != nil {
		t.Fatalf("prepare policy: %v", err)
	}
	return prepared
}

func evalPolicy(t *testing.T, prepared rego.PreparedEvalQuery, input AuthzInput) *AuthzResult {
	t.Helper()
	inputMap, err := structToMap(input)
	if err != nil {
		t.Fatalf("convert input: %v", err)
	}
	results, err := prepared.Eval(context.Background(), rego.EvalInput(inputMap))
	if err != nil {
		t.Fatalf("evaluate policy: %v", err)
	}
	decision, err := extractDecision(results)
	if err != nil {
		t.Fatalf("extract decision: %v", err)
	}
	return decision
}

func TestPublicEndpoint(t *testing.T) {
	prepared := prepareTestPolicy(t)

	result := evalPolicy(t, prepared, AuthzInput{
		Method: "/rootstock.v1.HealthService/Check",
	})

	if !result.Allow {
		t.Errorf("expected allow=true for public endpoint, got false")
	}
	if result.Reason != "public_endpoint" {
		t.Errorf("expected reason=public_endpoint, got %s", result.Reason)
	}
}

func TestPublicEndpointWithUser(t *testing.T) {
	prepared := prepareTestPolicy(t)

	result := evalPolicy(t, prepared, AuthzInput{
		SessionUserID: "user-123",
		Method:        "/rootstock.v1.HealthService/Check",
	})

	if !result.Allow {
		t.Errorf("expected allow=true for public endpoint with user, got false")
	}
	if result.Reason != "public_endpoint" {
		t.Errorf("expected reason=public_endpoint, got %s", result.Reason)
	}
}

func TestAuthenticatedNonPublic(t *testing.T) {
	prepared := prepareTestPolicy(t)

	result := evalPolicy(t, prepared, AuthzInput{
		SessionUserID: "user-123",
		Method:        "/rootstock.v1.SomeService/DoStuff",
	})

	if !result.Allow {
		t.Errorf("expected allow=true for authenticated user, got false")
	}
	if result.Reason != "authenticated" {
		t.Errorf("expected reason=authenticated, got %s", result.Reason)
	}
}

func TestDeniedUnauthenticatedNonPublic(t *testing.T) {
	prepared := prepareTestPolicy(t)

	result := evalPolicy(t, prepared, AuthzInput{
		SessionUserID: "",
		Method:        "/rootstock.v1.SomeService/DoStuff",
	})

	if result.Allow {
		t.Errorf("expected allow=false for unauthenticated non-public, got true")
	}
	if result.Reason != "denied" {
		t.Errorf("expected reason=denied, got %s", result.Reason)
	}
}

func TestDecisionExtraction(t *testing.T) {
	tests := []struct {
		name   string
		input  AuthzInput
		allow  bool
		reason string
	}{
		{
			name:   "health check no auth",
			input:  AuthzInput{Method: "/rootstock.v1.HealthService/Check"},
			allow:  true,
			reason: "public_endpoint",
		},
		{
			name:   "unknown method no auth",
			input:  AuthzInput{Method: "/rootstock.v1.Unknown/Method"},
			allow:  false,
			reason: "denied",
		},
		{
			name:   "unknown method with auth",
			input:  AuthzInput{SessionUserID: "u1", Method: "/rootstock.v1.Unknown/Method"},
			allow:  true,
			reason: "authenticated",
		},
	}

	prepared := prepareTestPolicy(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalPolicy(t, prepared, tt.input)
			if result.Allow != tt.allow {
				t.Errorf("allow: got %v, want %v", result.Allow, tt.allow)
			}
			if result.Reason != tt.reason {
				t.Errorf("reason: got %s, want %s", result.Reason, tt.reason)
			}
		})
	}
}
