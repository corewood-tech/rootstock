package server

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	"rootstock/web-server/auth"
	"rootstock/web-server/repo/authorization"
)

// AuthorizationInterceptor verifies the JWT (if present) and evaluates OPA policy.
// Decision logging happens inside the repository via OTel.
func AuthorizationInterceptor(verifier *JWTVerifier, authz authorization.Repository) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Extract and verify JWT from Authorization header
			authHeader := req.Header().Get("Authorization")
			subject, err := verifier.VerifyToken(ctx, authHeader)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("authentication failed: %w", err))
			}

			// Store subject in context for downstream handlers
			ctx = auth.ContextWithSubject(ctx, subject)

			input := authorization.AuthzInput{
				SessionUserID: subject,
				Method:        req.Spec().Procedure,
				Request:       req.Any(),
			}

			result, err := authz.Evaluate(ctx, input)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("authorization evaluation failed: %w", err))
			}

			if !result.Allow {
				return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("access denied: %s", result.Reason))
			}

			return next(ctx, req)
		}
	}
}

// BinaryOnlyInterceptor rejects requests that use JSON encoding,
// enforcing binary protobuf as the only accepted content type.
func BinaryOnlyInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ct := req.Header().Get("Content-Type")
			if strings.Contains(ct, "json") {
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("JSON encoding is not supported; use binary protobuf"),
				)
			}
			return next(ctx, req)
		}
	}
}
