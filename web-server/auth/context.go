package auth

import "context"

type contextKey string

const subjectKey contextKey = "subject"

// ContextWithSubject stores the authenticated subject (IdP user ID) in the context.
func ContextWithSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subjectKey, subject)
}

// SubjectFromContext extracts the authenticated subject (IdP user ID) from context.
func SubjectFromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(subjectKey).(string)
	return s, ok
}
