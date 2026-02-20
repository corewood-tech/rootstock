package device

// GetDeviceInput is what callers send to GetDeviceFlow.
type GetDeviceInput struct {
	DeviceID string
}

// RevokeDeviceInput is what callers send to RevokeDeviceFlow.
type RevokeDeviceInput struct {
	DeviceID string
}

// ReinstateDeviceInput is what callers send to ReinstateDeviceFlow.
type ReinstateDeviceInput struct {
	DeviceID string
}

// RegisterDeviceInput is what callers send to RegisterDeviceFlow.
type RegisterDeviceInput struct {
	EnrollmentCode string
	CSR            []byte // DER-encoded PKCS#10
}

// RenewCertInput is what callers send to RenewCertFlow.
type RenewCertInput struct {
	DeviceID string // from mTLS cert CN
	CSR      []byte // DER-encoded PKCS#10
}

// EnrollInCampaignInput is what callers send to EnrollInCampaignFlow.
type EnrollInCampaignInput struct {
	DeviceID   string
	CampaignID string
}
