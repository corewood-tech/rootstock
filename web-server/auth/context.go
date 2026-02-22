package auth

import "context"

type contextKey string

const (
	subjectKey      contextKey = "subject"
	sessionIDKey    contextKey = "session_id"
	sessionTokenKey contextKey = "session_token"
)

// ContextWithSubject stores the authenticated subject (IdP user ID) in the context.
func ContextWithSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subjectKey, subject)
}

// SubjectFromContext extracts the authenticated subject (IdP user ID) from context.
func SubjectFromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(subjectKey).(string)
	return s, ok
}

// ContextWithSessionID stores the session ID in the context.
func ContextWithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

// SessionIDFromContext extracts the session ID from context.
func SessionIDFromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(sessionIDKey).(string)
	return s, ok
}

// ContextWithSessionToken stores the session token in the context.
func ContextWithSessionToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, sessionTokenKey, token)
}

// SessionTokenFromContext extracts the session token from context.
func SessionTokenFromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(sessionTokenKey).(string)
	return s, ok
}
