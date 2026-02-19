package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func setupJWKS(t *testing.T) (jwk.Key, *httptest.Server) {
	t.Helper()

	raw, err := jwk.Import([]byte("test-secret-key-32-bytes-long!!!" ))
	if err != nil {
		t.Fatalf("create jwk from raw: %v", err)
	}
	if err := raw.Set(jwk.KeyIDKey, "test-key-id"); err != nil {
		t.Fatalf("set key id: %v", err)
	}
	if err := raw.Set(jwk.AlgorithmKey, jwa.HS256()); err != nil {
		t.Fatalf("set algorithm: %v", err)
	}

	set := jwk.NewSet()
	if err := set.AddKey(raw); err != nil {
		t.Fatalf("add key to set: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := json.Marshal(set)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}))
	t.Cleanup(srv.Close)

	return raw, srv
}

func signToken(t *testing.T, key jwk.Key, issuer, subject string) string {
	t.Helper()

	token, err := jwt.NewBuilder().
		Issuer(issuer).
		Subject(subject).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(time.Hour)).
		Build()
	if err != nil {
		t.Fatalf("build token: %v", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), key))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return string(signed)
}

func TestVerifyTokenValid(t *testing.T) {
	key, srv := setupJWKS(t)
	ctx := context.Background()
	issuer := "test-issuer"

	verifier, err := NewJWTVerifier(ctx, srv.URL, "localhost", issuer)
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	tokenStr := signToken(t, key, issuer, "user-123")
	subject, err := verifier.VerifyToken(ctx, "Bearer "+tokenStr)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if subject != "user-123" {
		t.Fatalf("expected subject user-123, got %s", subject)
	}
}

func TestVerifyTokenEmptyHeader(t *testing.T) {
	_, srv := setupJWKS(t)
	ctx := context.Background()

	verifier, err := NewJWTVerifier(ctx, srv.URL, "localhost", "test-issuer")
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	subject, err := verifier.VerifyToken(ctx, "")
	if err != nil {
		t.Fatalf("expected no error for empty header, got: %v", err)
	}
	if subject != "" {
		t.Fatalf("expected empty subject, got %s", subject)
	}
}

func TestVerifyTokenNoBearerPrefix(t *testing.T) {
	_, srv := setupJWKS(t)
	ctx := context.Background()

	verifier, err := NewJWTVerifier(ctx, srv.URL, "localhost", "test-issuer")
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	subject, err := verifier.VerifyToken(ctx, "Basic abc123")
	if err != nil {
		t.Fatalf("expected no error for non-bearer, got: %v", err)
	}
	if subject != "" {
		t.Fatalf("expected empty subject, got %s", subject)
	}
}

func TestVerifyTokenWrongIssuer(t *testing.T) {
	key, srv := setupJWKS(t)
	ctx := context.Background()

	verifier, err := NewJWTVerifier(ctx, srv.URL, "localhost", "correct-issuer")
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	tokenStr := signToken(t, key, "wrong-issuer", "user-456")
	_, err = verifier.VerifyToken(ctx, "Bearer "+tokenStr)
	if err == nil {
		t.Fatal("expected error for wrong issuer, got nil")
	}
}

func TestVerifyTokenExpired(t *testing.T) {
	key, srv := setupJWKS(t)
	ctx := context.Background()
	issuer := "test-issuer"

	verifier, err := NewJWTVerifier(ctx, srv.URL, "localhost", issuer)
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	// Build an already-expired token
	token, err := jwt.NewBuilder().
		Issuer(issuer).
		Subject("user-789").
		IssuedAt(time.Now().Add(-2 * time.Hour)).
		Expiration(time.Now().Add(-1 * time.Hour)).
		Build()
	if err != nil {
		t.Fatalf("build token: %v", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), key))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	_, err = verifier.VerifyToken(ctx, "Bearer "+string(signed))
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestVerifyTokenInvalidSignature(t *testing.T) {
	_, srv := setupJWKS(t)
	ctx := context.Background()
	issuer := "test-issuer"

	verifier, err := NewJWTVerifier(ctx, srv.URL, "localhost", issuer)
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	// Create a token signed with a different key
	differentKey, err := jwk.Import([]byte("different-secret-key-32-bytes!!!"))
	if err != nil {
		t.Fatalf("create different key: %v", err)
	}
	if err := differentKey.Set(jwk.KeyIDKey, "different-key"); err != nil {
		t.Fatalf("set key id: %v", err)
	}

	tokenStr := signToken(t, differentKey, issuer, "user-999")
	_, err = verifier.VerifyToken(ctx, "Bearer "+tokenStr)
	if err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}
}
