package cert

import "time"

// IssuedCert is the certificate record returned after issuance.
type IssuedCert struct {
	CertPEM   []byte
	Serial    string // hex-encoded
	NotBefore time.Time
	NotAfter  time.Time
}

// CACert is the CA certificate returned by GetCACert.
type CACert struct {
	CertPEM []byte
}
