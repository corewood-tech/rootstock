package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"testing"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"rootstock/web-server/config"
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

// makeDeviceCert creates a device cert signed by the test CA with specific timing.
func (ca *testCA) makeDeviceCert(t *testing.T, deviceID string, notBefore, notAfter time.Time) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()
	devKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate device key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: mustSerial(t),
		Subject:      pkix.Name{CommonName: deviceID},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	devDER, err := x509.CreateCertificate(rand.Reader, template, ca.CACert, &devKey.PublicKey, ca.CAKey)
	if err != nil {
		t.Fatalf("create device cert: %v", err)
	}
	devCert, _ := x509.ParseCertificate(devDER)
	return devCert, devKey
}

// makeTLSClient creates a mochi.Client with a real *tls.Conn containing specific peer certificates.
func makeTLSClient(t *testing.T, ca *testCA, deviceID string, deviceCert *x509.Certificate, deviceKey *ecdsa.PrivateKey) *mochi.Client {
	t.Helper()

	// Create TLS certificate pair for the device
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: deviceCert.Raw})
	keyDER, _ := x509.MarshalECPrivateKey(deviceKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	clientCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("x509 key pair: %v", err)
	}

	// Create server cert for the TLS server side
	serverKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	serverTemplate := &x509.Certificate{
		SerialNumber: mustSerial(t),
		Subject:      pkix.Name{CommonName: "localhost"},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	serverDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, ca.CACert, &serverKey.PublicKey, ca.CAKey)
	if err != nil {
		t.Fatalf("create server cert: %v", err)
	}
	serverCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER})
	serverKeyDER, _ := x509.MarshalECPrivateKey(serverKey)
	serverKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: serverKeyDER})
	serverTLSCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		t.Fatalf("server key pair: %v", err)
	}

	// Create an in-memory TLS connection pair
	serverConf := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		ClientCAs:    ca.CACertPool,
		ClientAuth:   tls.RequireAnyClientCert,
	}
	clientConf := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            ca.CACertPool,
		ServerName:         "localhost",
		InsecureSkipVerify: true,
	}

	serverConn, clientConn := net.Pipe()
	serverTLS := tls.Server(serverConn, serverConf)
	clientTLS := tls.Client(clientConn, clientConf)

	errCh := make(chan error, 2)
	go func() { errCh <- serverTLS.Handshake() }()
	go func() { errCh <- clientTLS.Handshake() }()

	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			serverTLS.Close()
			clientTLS.Close()
			t.Fatalf("tls handshake: %v", err)
		}
	}

	t.Cleanup(func() {
		serverTLS.Close()
		clientTLS.Close()
	})

	// The server-side connection has the peer certificates from the client
	return &mochi.Client{
		ID:  deviceID,
		Net: mochi.ClientConnection{Conn: serverTLS},
	}
}

func newHookWithGrace(t *testing.T, ca *testCA, graceDays int) *MQTTAuthHook {
	t.Helper()
	h := &MQTTAuthHook{}
	h.SetOpts(slog.Default(), &mochi.HookOptions{})
	if err := h.Init(&MQTTAuthHookConfig{CACertPool: ca.CACertPool, GracePeriodDays: graceDays}); err != nil {
		t.Fatalf("init hook: %v", err)
	}
	return h
}

func TestAuthHook_GracePeriodAllowsRenewAndCert(t *testing.T) {
	ca := newTestCA(t)
	h := newHookWithGrace(t, ca, 7)

	// Create a device cert that expired 2 days ago (within 7-day grace)
	now := time.Now()
	devCert, devKey := ca.makeDeviceCert(t, "device-grace",
		now.AddDate(0, -3, 0),  // notBefore: 3 months ago
		now.Add(-2*24*time.Hour), // notAfter: 2 days ago
	)

	cl := makeTLSClient(t, ca, "device-grace", devCert, devKey)

	// Should allow renew topic
	if !h.OnACLCheck(cl, "rootstock/device-grace/renew", true) {
		t.Error("grace period should allow renew topic")
	}

	// Should allow cert topic
	if !h.OnACLCheck(cl, "rootstock/device-grace/cert", false) {
		t.Error("grace period should allow cert topic")
	}
}

func TestAuthHook_GracePeriodDeniesDataTopics(t *testing.T) {
	ca := newTestCA(t)
	h := newHookWithGrace(t, ca, 7)

	now := time.Now()
	devCert, devKey := ca.makeDeviceCert(t, "device-grace2",
		now.AddDate(0, -3, 0),
		now.Add(-2*24*time.Hour),
	)

	cl := makeTLSClient(t, ca, "device-grace2", devCert, devKey)

	// Should deny data topics
	if h.OnACLCheck(cl, "rootstock/device-grace2/data/campaign-1", true) {
		t.Error("grace period should deny data topics")
	}

	// Should deny config topic
	if h.OnACLCheck(cl, "rootstock/device-grace2/config", false) {
		t.Error("grace period should deny config topics")
	}
}

func TestAuthHook_ValidCertAllowsAllTopics(t *testing.T) {
	ca := newTestCA(t)
	h := newHookWithGrace(t, ca, 7)

	now := time.Now()
	devCert, devKey := ca.makeDeviceCert(t, "device-valid",
		now.Add(-time.Hour),
		now.AddDate(0, 3, 0), // expires in 3 months — valid
	)

	cl := makeTLSClient(t, ca, "device-valid", devCert, devKey)

	topics := []struct {
		topic string
		write bool
	}{
		{"rootstock/device-valid/data/campaign-1", true},
		{"rootstock/device-valid/config", false},
		{"rootstock/device-valid/renew", true},
		{"rootstock/device-valid/cert", false},
	}

	for _, tc := range topics {
		if !h.OnACLCheck(cl, tc.topic, tc.write) {
			t.Errorf("valid cert should allow %s (write=%v)", tc.topic, tc.write)
		}
	}
}

func TestAuthHook_AuthenticateWithGracePeriodCert(t *testing.T) {
	ca := newTestCA(t)
	h := newHookWithGrace(t, ca, 7)

	now := time.Now()
	devCert, devKey := ca.makeDeviceCert(t, "device-auth-grace",
		now.AddDate(0, -3, 0),
		now.Add(-2*24*time.Hour), // expired 2 days ago
	)

	cl := makeTLSClient(t, ca, "device-auth-grace", devCert, devKey)

	// Should authenticate (TLS handshake already passed, hook just checks CN match)
	if !h.OnConnectAuthenticate(cl, packets.Packet{}) {
		t.Error("grace period device should authenticate (CN matches client ID)")
	}
}

func TestMQTTServer_GracePeriodCertAccepted(t *testing.T) {
	dir := t.TempDir()
	caCert, caKey := writeTestCA(t, dir)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	lis.Close()

	cfg := &config.Config{
		Cert: config.CertConfig{
			CACertPath: fmt.Sprintf("%s/ca.crt", dir),
			CAKeyPath:  fmt.Sprintf("%s/ca.key", dir),
		},
		MQTT: config.MQTTConfig{
			Port:            port,
			ServerSANs:      []string{"localhost"},
			GracePeriodDays: 7,
		},
	}

	mqttServer, cleanup, err := NewMQTTServer(cfg)
	if err != nil {
		t.Fatalf("NewMQTTServer(): %v", err)
	}
	defer cleanup()

	go mqttServer.Serve()
	time.Sleep(200 * time.Millisecond)

	// Create a device cert expired 2 days ago (within 7-day grace)
	now := time.Now()
	devKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	devTemplate := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "grace-device"},
		NotBefore:    now.AddDate(0, -3, 0),
		NotAfter:     now.Add(-2 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	devDER, _ := x509.CreateCertificate(rand.Reader, devTemplate, caCert, &devKey.PublicKey, caKey)
	devCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: devDER})
	devKeyDER, _ := x509.MarshalECPrivateKey(devKey)
	devKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: devKeyDER})
	deviceCert, _ := tls.X509KeyPair(devCertPEM, devKeyPEM)

	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{deviceCert},
		RootCAs:      caPool,
		ServerName:   "localhost",
	}

	// Grace period cert should be accepted at TLS level
	conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port), tlsCfg)
	if err != nil {
		t.Fatalf("TLS dial with grace-period cert should succeed: %v", err)
	}
	conn.Close()
}
