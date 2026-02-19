package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// JWTVerifier verifies JWTs against a JWKS endpoint.
type JWTVerifier struct {
	cache   *jwk.Cache
	jwksURL string
	issuer  string
}

// hostRewriteTransport overrides the Host header on outgoing requests.
// Zitadel resolves instances by hostname — internal Docker requests
// arrive with Host "zitadel:8080" but the instance is registered
// under the external domain (e.g. "localhost").
type hostRewriteTransport struct {
	host string
	base http.RoundTripper
}

func (t *hostRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = t.host
	return t.base.RoundTrip(req)
}

// NewJWTVerifier creates a verifier that fetches keys from the given JWKS URL.
// The hostOverride parameter sets the Host header sent to the JWKS endpoint,
// which is required when Zitadel's external domain differs from the container hostname.
func NewJWTVerifier(ctx context.Context, jwksURL string, hostOverride string, issuer string) (*JWTVerifier, error) {
	httpClient := &http.Client{
		Transport: &hostRewriteTransport{
			host: hostOverride,
			base: http.DefaultTransport,
		},
	}

	cache, err := jwk.NewCache(ctx, httprc.NewClient(httprc.WithHTTPClient(httpClient)))
	if err != nil {
		return nil, fmt.Errorf("create jwk cache: %w", err)
	}

	if err := cache.Register(ctx, jwksURL); err != nil {
		return nil, fmt.Errorf("register jwks url: %w", err)
	}

	return &JWTVerifier{
		cache:   cache,
		jwksURL: jwksURL,
		issuer:  issuer,
	}, nil
}

// VerifyToken parses and verifies a JWT, returning the subject claim.
// Returns empty string if the token is missing or invalid.
func (v *JWTVerifier) VerifyToken(ctx context.Context, authHeader string) (string, error) {
	if authHeader == "" {
		return "", nil
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return "", nil // no Bearer prefix
	}

	keyset, err := v.cache.Lookup(ctx, v.jwksURL)
	if err != nil {
		return "", fmt.Errorf("fetch jwks: %w", err)
	}

	token, err := jwt.Parse([]byte(tokenStr),
		jwt.WithKeySet(keyset),
		jwt.WithIssuer(v.issuer),
	)
	if err != nil {
		return "", fmt.Errorf("parse/verify jwt: %w", err)
	}

	subject, _ := token.Subject()
	return subject, nil
}
