package server

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	"rootstock/web-server/auth"
	userops "rootstock/web-server/ops/user"
	"rootstock/web-server/repo/authorization"
)

// AuthorizationInterceptor verifies the session (if present), resolves principal
// attributes, and evaluates OPA ABAC policy.
// Authorization header format: "Bearer sessionID|sessionToken"
func AuthorizationInterceptor(uOps *userops.Ops, authz authorization.Repository) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			var subject string
			var principal authorization.Principal

			authHeader := req.Header().Get("Authorization")
			if authHeader != "" {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				if tokenStr != authHeader {
					parts := strings.SplitN(tokenStr, "|", 2)
					if len(parts) != 2 {
						return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid authorization format"))
					}

					sessionID, sessionToken := parts[0], parts[1]
					validated, err := uOps.ValidateSession(ctx, userops.ValidateSessionInput{
						SessionID:    sessionID,
						SessionToken: sessionToken,
					})
					if err != nil {
						return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("session validation failed: %w", err))
					}

					subject = validated.UserID
					ctx = auth.ContextWithSessionID(ctx, sessionID)
					ctx = auth.ContextWithSessionToken(ctx, sessionToken)

					// Resolve principal attributes from app DB.
					appUser, err := uOps.GetUserByIdpID(ctx, subject)
					if err == nil && appUser != nil {
						principal.Role = appUser.UserType
					}
				}
			}

			ctx = auth.ContextWithSubject(ctx, subject)

			input := authorization.AuthzInput{
				SessionUserID: subject,
				Principal:     principal,
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
