package auth

// User type constants matching the DB CHECK constraint.
const (
	RoleResearcher = "researcher"
	RoleScitizen   = "scitizen"
	RoleBoth       = "both"
)

// IsValidRegistrationRole returns true for roles allowed at registration.
func IsValidRegistrationRole(role string) bool {
	return role == RoleResearcher || role == RoleScitizen
}

// IsValidUserType returns true for any valid user_type value.
func IsValidUserType(role string) bool {
	return role == RoleResearcher || role == RoleScitizen || role == RoleBoth
}
