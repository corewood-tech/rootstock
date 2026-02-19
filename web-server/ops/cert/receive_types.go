package cert

// IssueCertInput is what callers send to IssueCert.
type IssueCertInput struct {
	DeviceID string
	CSR      []byte // DER-encoded PKCS#10
}
