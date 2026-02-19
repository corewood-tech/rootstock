package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	deviceflows "rootstock/web-server/flows/device"
	httphandlers "rootstock/web-server/handlers/http"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	certrepo "rootstock/web-server/repo/cert"
	devicerepo "rootstock/web-server/repo/device"
)

// NewIoTServer wires cert repo → cert ops → device flows → HTTP handlers
// and returns an http.Handler + TLS config + a shutdown function.
func NewIoTServer(cfg *config.Config, pool *pgxpool.Pool) (http.Handler, *tls.Config, func(), error) {
	// Cert repo (in-process CA)
	crtRepo, err := certrepo.NewRepository(cfg.Cert)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create cert repo: %w", err)
	}

	// Device repo
	dRepo := devicerepo.NewRepository(pool)

	// Ops
	crtOps := certops.NewOps(crtRepo)
	dOps := deviceops.NewOps(dRepo)

	// Flows
	registerDeviceFlow := deviceflows.NewRegisterDeviceFlow(dOps, crtOps)
	renewCertFlow := deviceflows.NewRenewCertFlow(dOps, crtOps)

	// HTTP handlers
	deviceHandler := httphandlers.NewDeviceHandler(registerDeviceFlow, renewCertFlow, crtOps)

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

	// Load server certificate (same CA cert+key used as server identity for device port)
	serverCert, err := tls.LoadX509KeyPair(cfg.Cert.CACertPath, cfg.Cert.CAKeyPath)
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
		dRepo.Shutdown()
	}

	return mux, tlsConfig, shutdown, nil
}
