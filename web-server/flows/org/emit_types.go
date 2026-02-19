package org

// Org is the organization record returned by org flows.
type Org struct {
	ID   string
	Name string
}

// Role is the role record returned by DefineRoleFlow.
type Role struct {
	ProjectID   string
	RoleKey     string
	DisplayName string
}

// UserGrant is the grant record returned by AssignRoleFlow.
type UserGrant struct {
	UserGrantID string
	UserID      string
	ProjectID   string
	RoleKeys    []string
}

// InviteResult is the result returned by InviteUserFlow.
type InviteResult struct {
	UserID    string
	EmailCode string
}

// OnboardInstitutionResult is the result returned by OnboardInstitutionFlow.
type OnboardInstitutionResult struct {
	OrgID  string
	UserID string
}
