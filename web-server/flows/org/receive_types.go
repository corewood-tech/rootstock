package org

// CreateOrgInput is what callers send to CreateOrgFlow.
type CreateOrgInput struct {
	Name string
}

// NestOrgInput is what callers send to NestOrgFlow.
type NestOrgInput struct {
	Name        string
	ParentOrgID string
}

// DefineRoleInput is what callers send to DefineRoleFlow.
type DefineRoleInput struct {
	ProjectID   string
	RoleKey     string
	DisplayName string
}

// AssignRoleInput is what callers send to AssignRoleFlow.
type AssignRoleInput struct {
	UserID    string
	ProjectID string
	RoleKeys  []string
}

// InviteUserInput is what callers send to InviteUserFlow.
type InviteUserInput struct {
	OrgID      string
	Email      string
	GivenName  string
	FamilyName string
}

// OnboardInstitutionInput is what callers send to OnboardInstitutionFlow.
type OnboardInstitutionInput struct {
	OrgName         string
	ParentOrgID     string
	ProjectID       string
	AdminEmail      string
	AdminGivenName  string
	AdminFamilyName string
}
