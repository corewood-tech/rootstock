package http

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	certops "rootstock/web-server/ops/cert"
	deviceflows "rootstock/web-server/flows/device"
)

// DeviceHandler serves device enrollment and certificate endpoints.
type DeviceHandler struct {
	registerDevice *deviceflows.RegisterDeviceFlow
	renewCert      *deviceflows.RenewCertFlow
	certOps        *certops.Ops
}

// NewDeviceHandler creates the handler with all required flows.
func NewDeviceHandler(
	registerDevice *deviceflows.RegisterDeviceFlow,
	renewCert *deviceflows.RenewCertFlow,
	certOps *certops.Ops,
) *DeviceHandler {
	return &DeviceHandler{
		registerDevice: registerDevice,
		renewCert:      renewCert,
		certOps:        certOps,
	}
}

type enrollRequest struct {
	EnrollmentCode string `json:"enrollment_code"`
	CSR            string `json:"csr"` // PEM-encoded CSR
}

type enrollResponse struct {
	DeviceID  string `json:"device_id"`
	CertPEM   string `json:"cert_pem"`
	Serial    string `json:"serial"`
	NotBefore string `json:"not_before"`
	NotAfter  string `json:"not_after"`
}

type renewRequest struct {
	CSR string `json:"csr"` // PEM-encoded CSR
}

type renewResponse struct {
	CertPEM   string `json:"cert_pem"`
	Serial    string `json:"serial"`
	NotBefore string `json:"not_before"`
	NotAfter  string `json:"not_after"`
}

// Enroll handles POST /enroll — no mTLS required, auth is via enrollment code.
func (h *DeviceHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req enrollRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.EnrollmentCode == "" || req.CSR == "" {
		http.Error(w, "enrollment_code and csr are required", http.StatusBadRequest)
		return
	}

	csrDER, err := decodePEMToCSR([]byte(req.CSR))
	if err != nil {
		http.Error(w, "invalid csr: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.registerDevice.Run(r.Context(), deviceflows.RegisterDeviceInput{
		EnrollmentCode: req.EnrollmentCode,
		CSR:            csrDER,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := enrollResponse{
		DeviceID:  result.DeviceID,
		CertPEM:   string(result.CertPEM),
		Serial:    result.Serial,
		NotBefore: result.NotBefore.Format(time.RFC3339),
		NotAfter:  result.NotAfter.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Renew handles POST /renew — mTLS required, device ID from cert CN.
func (h *DeviceHandler) Renew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract device ID from mTLS peer certificate
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		http.Error(w, "client certificate required", http.StatusUnauthorized)
		return
	}

	peerCert := r.TLS.PeerCertificates[0]
	deviceID := peerCert.Subject.CommonName

	// Grace period: accept certs expired <= 7 days
	if time.Now().After(peerCert.NotAfter.Add(7 * 24 * time.Hour)) {
		http.Error(w, "client certificate expired beyond grace period", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req renewRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.CSR == "" {
		http.Error(w, "csr is required", http.StatusBadRequest)
		return
	}

	csrDER, err := decodePEMToCSR([]byte(req.CSR))
	if err != nil {
		http.Error(w, "invalid csr: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.renewCert.Run(r.Context(), deviceflows.RenewCertInput{
		DeviceID: deviceID,
		CSR:      csrDER,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := renewResponse{
		CertPEM:   string(result.CertPEM),
		Serial:    result.Serial,
		NotBefore: result.NotBefore.Format(time.RFC3339),
		NotAfter:  result.NotAfter.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCACert handles GET /ca — public, returns CA cert PEM.
func (h *DeviceHandler) GetCACert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ca, err := h.certOps.GetCACert(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Write(ca.CertPEM)
}

func decodePEMToCSR(data []byte) ([]byte, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}
	if block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("expected CERTIFICATE REQUEST, got %s", block.Type)
	}
	return block.Bytes, nil
}
