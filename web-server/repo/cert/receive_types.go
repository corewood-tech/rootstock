package cert

// IssueCertInput is what the IssueCert op sends to the repository.
type IssueCertInput struct {
	DeviceID string // becomes certificate CN — caller determines identity, not the CSR
	CSR      []byte // DER-encoded PKCS#10
}
