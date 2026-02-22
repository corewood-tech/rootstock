package org

// Org is the organization record returned by org ops.
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

// CreatedIdpUser is the result of creating a user in the IdP.
type CreatedIdpUser struct {
	UserID string
}
