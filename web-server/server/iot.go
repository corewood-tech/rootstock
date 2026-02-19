package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"rootstock/web-server/config"
	deviceflows "rootstock/web-server/flows/device"
	httphandlers "rootstock/web-server/handlers/http"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	certrepo "rootstock/web-server/repo/cert"
)

// NewIoTServer wires cert repo → cert ops → device flows → HTTP handlers
// and returns an http.Handler + TLS config + a shutdown function.
// It takes shared device ops to avoid duplicating the device repo.
func NewIoTServer(cfg *config.Config, dOps *deviceops.Ops) (http.Handler, *tls.Config, func(), error) {
	// Cert repo (in-process CA)
	crtRepo, err := certrepo.NewRepository(cfg.Cert)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create cert repo: %w", err)
	}

	// Ops
	crtOps := certops.NewOps(crtRepo)

	// Flows
	registerDeviceFlow := deviceflows.NewRegisterDeviceFlow(dOps, crtOps)
	renewCertFlow := deviceflows.NewRenewCertFlow(dOps, crtOps)
	getCACertFlow := deviceflows.NewGetCACertFlow(crtOps)

	// HTTP handlers
	deviceHandler := httphandlers.NewDeviceHandler(registerDeviceFlow, renewCertFlow, getCACertFlow)

	mux := http.NewServeMux()
	mux.HandleFunc("/enroll", deviceHandler.Enroll)
	mux.HandleFunc("/renew", deviceHandler.Renew)
	mux.HandleFunc("/ca", deviceHandler.GetCACert)

	// TLS config with optional client certs
	caCertPEM, err := os.ReadFile(cfg.Cert.CACertPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read ca cert for tls: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCertPEM) {
		return nil, nil, nil, fmt.Errorf("failed to add ca cert to pool")
	}

	// Load server leaf certificate (issued by the CA, with serverAuth EKU)
	serverCert, err := tls.LoadX509KeyPair(cfg.Cert.ServerCertPath, cfg.Cert.ServerKeyPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load server cert: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.VerifyClientCertIfGiven, // /enroll has no cert, /renew does
	}

	shutdown := func() {
		crtRepo.Shutdown()
	}

	return mux, tlsConfig, shutdown, nil
}
