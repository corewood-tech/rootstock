package cert

import "context"

// Repository defines the interface for certificate authority operations.
type Repository interface {
	IssueCert(ctx context.Context, input IssueCertInput) (*IssuedCert, error)
	GetCACert(ctx context.Context) (*CACert, error)
	Shutdown()
}
