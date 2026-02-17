package authorization

import "context"

// Repository defines the interface for authorization operations.
type Repository interface {
	// Evaluate checks if the given input is authorized.
	Evaluate(ctx context.Context, input AuthzInput) (*AuthzResult, error)

	// Recompile regenerates the policy from current state.
	Recompile(ctx context.Context) error
}
