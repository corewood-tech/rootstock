package device

import (
	"context"

	certops "rootstock/web-server/ops/cert"
)

// GetCACertFlow retrieves the CA certificate.
type GetCACertFlow struct {
	certOps *certops.Ops
}

// NewGetCACertFlow creates the flow with its required ops.
func NewGetCACertFlow(certOps *certops.Ops) *GetCACertFlow {
	return &GetCACertFlow{certOps: certOps}
}

// Run returns the CA certificate PEM.
func (f *GetCACertFlow) Run(ctx context.Context) (*CACert, error) {
	result, err := f.certOps.GetCACert(ctx)
	if err != nil {
		return nil, err
	}
	return &CACert{CertPEM: result.CertPEM}, nil
}
