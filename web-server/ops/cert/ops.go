package cert

import (
	"context"

	certrepo "rootstock/web-server/repo/cert"
)

// Ops holds certificate operations. Each method is one op.
type Ops struct {
	repo certrepo.Repository
}

// NewOps creates cert ops backed by the given repository.
func NewOps(repo certrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// IssueCert issues a client certificate for a device.
func (o *Ops) IssueCert(ctx context.Context, input IssueCertInput) (*IssuedCert, error) {
	result, err := o.repo.IssueCert(ctx, toRepoIssueCertInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoIssuedCert(result), nil
}

// GetCACert returns the CA certificate PEM.
func (o *Ops) GetCACert(ctx context.Context) (*CACert, error) {
	result, err := o.repo.GetCACert(ctx)
	if err != nil {
		return nil, err
	}
	return &CACert{CertPEM: result.CertPEM}, nil
}

func toRepoIssueCertInput(in IssueCertInput) certrepo.IssueCertInput {
	return certrepo.IssueCertInput{
		DeviceID: in.DeviceID,
		CSR:      in.CSR,
	}
}

func fromRepoIssuedCert(r *certrepo.IssuedCert) *IssuedCert {
	return &IssuedCert{
		CertPEM:   r.CertPEM,
		Serial:    r.Serial,
		NotBefore: r.NotBefore,
		NotAfter:  r.NotAfter,
	}
}
