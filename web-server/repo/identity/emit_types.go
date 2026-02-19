package identity

// Org is the organization record returned by identity operations.
type Org struct {
	ID   string
	Name string
}

// Role is the role record returned by DefineRole.
type Role struct {
	ProjectID   string
	RoleKey     string
	DisplayName string
}

// UserGrant is the grant record returned by AssignRole.
type UserGrant struct {
	UserGrantID string
	UserID      string
	ProjectID   string
	RoleKeys    []string
}

// InviteResult is the result of inviting a user.
type InviteResult struct {
	UserID    string
	EmailCode string
}
