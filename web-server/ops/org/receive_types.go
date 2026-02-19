package org

// CreateOrgInput is what callers send to CreateOrg.
type CreateOrgInput struct {
	Name string
}

// NestOrgInput is what callers send to NestOrg.
type NestOrgInput struct {
	Name        string
	ParentOrgID string
}

// DefineRoleInput is what callers send to DefineRole.
type DefineRoleInput struct {
	ProjectID   string
	RoleKey     string
	DisplayName string
}

// AssignRoleInput is what callers send to AssignRole.
type AssignRoleInput struct {
	UserID    string
	ProjectID string
	RoleKeys  []string
}

// InviteUserInput is what callers send to InviteUser.
type InviteUserInput struct {
	OrgID      string
	Email      string
	GivenName  string
	FamilyName string
}
