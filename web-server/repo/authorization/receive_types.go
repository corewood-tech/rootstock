package authorization

// Principal describes the authenticated subject's attributes for ABAC evaluation.
type Principal struct {
	Role string // "researcher", "scitizen", "both", or "" (unauthenticated)
}

// AuthzInput is what the interceptor sends to the authorization repository.
type AuthzInput struct {
	SessionUserID string
	Principal     Principal
	Method        string
	Request       interface{}
}

// ToOPAInput builds the map that OPA expects, without JSON serialization.
func (a AuthzInput) ToOPAInput() map[string]interface{} {
	return map[string]interface{}{
		"session_user_id": a.SessionUserID,
		"principal": map[string]interface{}{
			"role": a.Principal.Role,
		},
		"method":  a.Method,
		"request": a.Request,
	}
}
