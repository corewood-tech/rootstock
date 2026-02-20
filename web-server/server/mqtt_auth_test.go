package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"log/slog"
	"math/big"
	"testing"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

// testCA generates a throwaway CA + device cert for testing.
type testCA struct {
	CACert     *x509.Certificate
	CACertPool *x509.CertPool
	CAKey      *ecdsa.PrivateKey
}

func newTestCA(t *testing.T) *testCA {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ca key: %v", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber:          mustSerial(t),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create ca cert: %v", err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parse ca cert: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(caCert)

	return &testCA{CACert: caCert, CACertPool: pool, CAKey: caKey}
}


func mustSerial(t *testing.T) *big.Int {
	t.Helper()
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	s, err := rand.Int(rand.Reader, max)
	if err != nil {
		t.Fatalf("generate serial: %v", err)
	}
	return s
}

func newHook(t *testing.T, ca *testCA) *MQTTAuthHook {
	t.Helper()
	h := &MQTTAuthHook{}
	h.SetOpts(slog.Default(), &mochi.HookOptions{})
	if err := h.Init(&MQTTAuthHookConfig{CACertPool: ca.CACertPool}); err != nil {
		t.Fatalf("init hook: %v", err)
	}
	return h
}

func TestAuthHook_Provides(t *testing.T) {
	h := &MQTTAuthHook{}

	if !h.Provides(mochi.OnConnectAuthenticate) {
		t.Error("should provide OnConnectAuthenticate")
	}
	if !h.Provides(mochi.OnACLCheck) {
		t.Error("should provide OnACLCheck")
	}
	if h.Provides(mochi.OnPublish) {
		t.Error("should not provide OnPublish")
	}
}

func TestAuthHook_AuthenticateInlineClient(t *testing.T) {
	ca := newTestCA(t)
	h := newHook(t, ca)

	cl := &mochi.Client{
		Net: mochi.ClientConnection{Inline: true},
	}

	if !h.OnConnectAuthenticate(cl, packets.Packet{}) {
		t.Error("inline client should always authenticate")
	}
}

func TestAuthHook_RejectNonTLSConnection(t *testing.T) {
	ca := newTestCA(t)
	h := newHook(t, ca)

	// A plain net.Conn (not *tls.Conn) should be rejected
	cl := &mochi.Client{
		ID:  "device-001",
		Net: mochi.ClientConnection{},
	}

	if h.OnConnectAuthenticate(cl, packets.Packet{}) {
		t.Error("non-TLS connection should reject")
	}
}

func TestAuthHook_ACLAllowOwnTopic(t *testing.T) {
	ca := newTestCA(t)
	h := newHook(t, ca)

	cl := &mochi.Client{
		ID:  "device-001",
		Net: mochi.ClientConnection{},
	}

	cases := []struct {
		topic string
		write bool
	}{
		{"rootstock/device-001/data/campaign-1", true},
		{"rootstock/device-001/config", false},
		{"rootstock/device-001/renew", true},
		{"rootstock/device-001/cert", false},
	}

	for _, tc := range cases {
		if !h.OnACLCheck(cl, tc.topic, tc.write) {
			t.Errorf("should allow %s (write=%v) for own device", tc.topic, tc.write)
		}
	}
}

func TestAuthHook_ACLRejectCrossDevice(t *testing.T) {
	ca := newTestCA(t)
	h := newHook(t, ca)

	cl := &mochi.Client{
		ID:  "device-001",
		Net: mochi.ClientConnection{},
	}

	cases := []string{
		"rootstock/device-OTHER/data/campaign-1",
		"rootstock/device-OTHER/config",
	}

	for _, topic := range cases {
		if h.OnACLCheck(cl, topic, true) {
			t.Errorf("should reject cross-device access to %s", topic)
		}
	}
}

func TestAuthHook_ACLRejectInvalidPrefix(t *testing.T) {
	ca := newTestCA(t)
	h := newHook(t, ca)

	cl := &mochi.Client{
		ID:  "device-001",
		Net: mochi.ClientConnection{},
	}

	if h.OnACLCheck(cl, "other/device-001/data", true) {
		t.Error("should reject non-rootstock prefix")
	}
}

func TestAuthHook_ACLAllowInlineClient(t *testing.T) {
	ca := newTestCA(t)
	h := newHook(t, ca)

	cl := &mochi.Client{
		ID:  "inline",
		Net: mochi.ClientConnection{Inline: true},
	}

	if !h.OnACLCheck(cl, "rootstock/any-device/config", true) {
		t.Error("inline client should have unrestricted ACL access")
	}
}
