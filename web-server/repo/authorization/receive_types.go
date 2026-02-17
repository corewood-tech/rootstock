package authorization

// AuthzInput is what the interceptor sends to the authorization repository.
type AuthzInput struct {
	SessionUserID string      `json:"session_user_id"`
	Method        string      `json:"method"`
	Request       interface{} `json:"request"`
}
