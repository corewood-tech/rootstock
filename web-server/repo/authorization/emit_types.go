package authorization

// AuthzResult is what the authorization repository returns.
type AuthzResult struct {
	Allow  bool
	Reason string // "public_endpoint", "authenticated", "denied"
}
