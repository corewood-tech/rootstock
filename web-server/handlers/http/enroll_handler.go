package http

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	deviceflows "rootstock/web-server/flows/device"
)

// EnrollHandler serves device enrollment and CA cert endpoints on the RPC port.
// No mTLS — enrollment uses enrollment code auth, CA cert is public.
type EnrollHandler struct {
	registerDevice *deviceflows.RegisterDeviceFlow
	getCACert      *deviceflows.GetCACertFlow
}

// NewEnrollHandler creates the handler with required flows.
func NewEnrollHandler(
	registerDevice *deviceflows.RegisterDeviceFlow,
	getCACert *deviceflows.GetCACertFlow,
) *EnrollHandler {
	return &EnrollHandler{
		registerDevice: registerDevice,
		getCACert:      getCACert,
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

// Enroll handles POST /enroll — no mTLS required, auth is via enrollment code.
func (h *EnrollHandler) Enroll(w http.ResponseWriter, r *http.Request) {
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
		writeFlowError(w, err)
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

// GetCACert handles GET /ca — public, returns CA cert PEM.
func (h *EnrollHandler) GetCACert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ca, err := h.getCACert.Run(r.Context())
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

// writeFlowError maps domain errors to appropriate HTTP status codes.
func writeFlowError(w http.ResponseWriter, err error) {
	msg := err.Error()

	if strings.Contains(msg, "not found") ||
		strings.Contains(msg, "expired") ||
		strings.Contains(msg, "already used") {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	if strings.Contains(msg, "parse csr") ||
		strings.Contains(msg, "csr signature") ||
		strings.Contains(msg, "key too small") ||
		strings.Contains(msg, "unsupported key type") {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	http.Error(w, msg, http.StatusInternalServerError)
}
