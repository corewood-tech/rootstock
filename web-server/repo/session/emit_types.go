package session

// Session is the session record returned by session operations.
type Session struct {
	SessionID    string
	SessionToken string
	UserID       string
}
