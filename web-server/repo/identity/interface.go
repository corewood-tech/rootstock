package identity

import "context"

// Repository defines the interface for identity management operations.
type Repository interface {
	CreateOrg(ctx context.Context, input CreateOrgInput) (*Org, error)
	NestOrg(ctx context.Context, input NestOrgInput) (*Org, error)
	DefineRole(ctx context.Context, input DefineRoleInput) (*Role, error)
	AssignRole(ctx context.Context, input AssignRoleInput) (*UserGrant, error)
	InviteUser(ctx context.Context, input InviteUserInput) (*InviteResult, error)
	Shutdown()
}
