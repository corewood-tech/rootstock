package identity

// CreateOrgInput is what callers send to create an organization.
type CreateOrgInput struct {
	Name string
}

// NestOrgInput is what callers send to create a child organization.
type NestOrgInput struct {
	Name       string
	ParentOrgID string
}

// DefineRoleInput is what callers send to define a project role.
type DefineRoleInput struct {
	ProjectID   string
	RoleKey     string
	DisplayName string
}

// AssignRoleInput is what callers send to grant roles to a user.
type AssignRoleInput struct {
	UserID    string
	ProjectID string
	RoleKeys  []string
}

// InviteUserInput is what callers send to create and invite a user.
type InviteUserInput struct {
	OrgID      string
	Email      string
	GivenName  string
	FamilyName string
}

// CreateHumanUserInput is what callers send to create a user with a password.
type CreateHumanUserInput struct {
	Email      string
	Password   string
	GivenName  string
	FamilyName string
}

// VerifyEmailInput is what callers send to verify a user's email address.
type VerifyEmailInput struct {
	UserID           string
	VerificationCode string
}
