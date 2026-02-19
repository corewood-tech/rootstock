package cert

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"rootstock/web-server/config"
)

type response[T any] struct {
	val T
	err error
}

type issueReq struct {
	ctx   context.Context
	input IssueCertInput
	resp  chan response[*IssuedCert]
}

type getCACertReq struct {
	ctx  context.Context
	resp chan response[*CACert]
}

type shutdownReq struct {
	resp chan struct{}
}

type caRepo struct {
	caCert       *x509.Certificate
	caCertPEM    []byte
	caSigner     crypto.Signer
	lifetimeDays int
	issueCh      chan issueReq
	getCACertCh  chan getCACertReq
	shutdownCh   chan shutdownReq
}

// NewRepository creates a cert repository backed by an in-process X.509 CA.
func NewRepository(cfg config.CertConfig) (Repository, error) {
	// Load CA cert
	certPEM, err := os.ReadFile(cfg.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("read ca cert: %w", err)
	}
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("decode ca cert pem: no PEM block found")
	}
	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse ca cert: %w", err)
	}

	// Load CA key — skip EC PARAMETERS block if present
	keyPEM, err := os.ReadFile(cfg.CAKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read ca key: %w", err)
	}
	var keyBlock *pem.Block
	rest := keyPEM
	for {
		keyBlock, rest = pem.Decode(rest)
		if keyBlock == nil {
			return nil, fmt.Errorf("decode ca key pem: no EC PRIVATE KEY block found")
		}
		if keyBlock.Type == "EC PRIVATE KEY" {
			break
		}
	}
	caKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse ca key: %w", err)
	}

	// Assert key implements crypto.Signer
	signer, ok := any(caKey).(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("ca key does not implement crypto.Signer")
	}

	r := &caRepo{
		caCert:       caCert,
		caCertPEM:    certPEM,
		caSigner:     signer,
		lifetimeDays: cfg.CertLifetimeDays,
		issueCh:      make(chan issueReq),
		getCACertCh:  make(chan getCACertReq),
		shutdownCh:   make(chan shutdownReq),
	}
	go r.manage()
	return r, nil
}

func (r *caRepo) manage() {
	for {
		select {
		case req := <-r.issueCh:
			val, err := r.doIssue(req.ctx, req.input)
			req.resp <- response[*IssuedCert]{val: val, err: err}
		case req := <-r.getCACertCh:
			req.resp <- response[*CACert]{val: &CACert{CertPEM: r.caCertPEM}}
		case req := <-r.shutdownCh:
			close(req.resp)
			return
		}
	}
}

func (r *caRepo) IssueCert(ctx context.Context, input IssueCertInput) (*IssuedCert, error) {
	resp := make(chan response[*IssuedCert], 1)
	r.issueCh <- issueReq{ctx: ctx, input: input, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *caRepo) GetCACert(ctx context.Context) (*CACert, error) {
	resp := make(chan response[*CACert], 1)
	r.getCACertCh <- getCACertReq{ctx: ctx, resp: resp}
	res := <-resp
	return res.val, res.err
}

func (r *caRepo) Shutdown() {
	resp := make(chan struct{}, 1)
	r.shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

func (r *caRepo) doIssue(_ context.Context, input IssueCertInput) (*IssuedCert, error) {
	// Parse DER-encoded CSR
	csr, err := x509.ParseCertificateRequest(input.CSR)
	if err != nil {
		return nil, fmt.Errorf("parse csr: %w", err)
	}

	// Verify CSR is self-consistent
	if err := csr.CheckSignature(); err != nil {
		return nil, fmt.Errorf("csr signature check: %w", err)
	}

	// Validate key type
	switch pub := csr.PublicKey.(type) {
	case *ecdsa.PublicKey:
		if pub.Curve.Params().BitSize < 256 {
			return nil, fmt.Errorf("ecdsa key too small: %d bits", pub.Curve.Params().BitSize)
		}
	case *rsa.PublicKey:
		if pub.N.BitLen() < 2048 {
			return nil, fmt.Errorf("rsa key too small: %d bits", pub.N.BitLen())
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %T", csr.PublicKey)
	}

	// Generate 128-bit random serial
	serialMax := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialMax)
	if err != nil {
		return nil, fmt.Errorf("generate serial: %w", err)
	}

	now := time.Now().UTC()
	notAfter := now.Add(time.Duration(r.lifetimeDays) * 24 * time.Hour)

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: input.DeviceID},
		NotBefore:    now,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, r.caCert, csr.PublicKey, r.caSigner)
	if err != nil {
		return nil, fmt.Errorf("create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return &IssuedCert{
		CertPEM:   certPEM,
		Serial:    fmt.Sprintf("%x", serial),
		NotBefore: now,
		NotAfter:  notAfter,
	}, nil
}
