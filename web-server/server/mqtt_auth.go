package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"strings"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

// MQTTAuthHook implements mochi-mqtt's Hook interface for mTLS device
// authentication and topic-level ACL enforcement.
//
// Authentication: extracts device ID from the client certificate's CommonName.
// The TLS listener uses RequireAnyClientCert + custom VerifyPeerCertificate,
// so chain validation (with grace-period expiry) is already done at TLS level.
// This hook just extracts the identity.
//
// ACL: devices can only publish/subscribe to rootstock/{own-device-id}/*.
type MQTTAuthHook struct {
	mochi.HookBase
	caCertPool      *x509.CertPool
	gracePeriodDays int
}

// MQTTAuthHookConfig holds configuration for the auth hook.
type MQTTAuthHookConfig struct {
	CACertPool      *x509.CertPool
	GracePeriodDays int
}

func (h *MQTTAuthHook) ID() string {
	return "mqtt-mtls-auth"
}

func (h *MQTTAuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mochi.OnConnectAuthenticate,
		mochi.OnACLCheck,
	}, []byte{b})
}

func (h *MQTTAuthHook) Init(config any) error {
	if cfg, ok := config.(*MQTTAuthHookConfig); ok && cfg != nil {
		h.caCertPool = cfg.CACertPool
		h.gracePeriodDays = cfg.GracePeriodDays
	}
	return nil
}

// OnConnectAuthenticate verifies the client presented a valid mTLS certificate
// issued by our CA. The device ID is the certificate's CommonName, which must
// match the MQTT client ID.
func (h *MQTTAuthHook) OnConnectAuthenticate(cl *mochi.Client, pk packets.Packet) bool {
	// Allow the inline client (server-side publish/subscribe)
	if cl.Net.Inline {
		return true
	}

	tlsConn, ok := cl.Net.Conn.(*tls.Conn)
	if !ok {
		h.Log.Warn("mqtt auth: connection is not TLS", "client", cl.ID)
		return false
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		h.Log.Warn("mqtt auth: no peer certificate", "client", cl.ID)
		return false
	}

	peerCert := state.PeerCertificates[0]
	deviceID := peerCert.Subject.CommonName

	// The MQTT client ID must match the certificate's CN
	if cl.ID != deviceID {
		h.Log.Warn("mqtt auth: client ID mismatch",
			"client_id", cl.ID,
			"cert_cn", deviceID)
		return false
	}

	h.Log.Info("mqtt auth: device authenticated",
		"device_id", deviceID,
		"serial", peerCert.SerialNumber.String())
	return true
}

// isInGracePeriod re-derives grace status from the TLS connection state.
// Returns true if the client's cert is expired but within the grace window.
func (h *MQTTAuthHook) isInGracePeriod(cl *mochi.Client) bool {
	tlsConn, ok := cl.Net.Conn.(*tls.Conn)
	if !ok {
		return false
	}
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return false
	}
	return time.Now().After(state.PeerCertificates[0].NotAfter)
}

// OnACLCheck enforces that devices can only access their own topic namespace:
// rootstock/{device-id}/*. The device ID comes from the MQTT client ID,
// which was verified against the cert CN in OnConnectAuthenticate.
//
// Grace period: devices with expired (but within grace window) certs can only
// access renew and cert subtopics.
func (h *MQTTAuthHook) OnACLCheck(cl *mochi.Client, topic string, write bool) bool {
	// Allow the inline client (server-side operations)
	if cl.Net.Inline {
		return true
	}

	// Topic format: rootstock/{device-id}/...
	parts := strings.SplitN(topic, "/", 3)
	if len(parts) < 2 || parts[0] != "rootstock" {
		h.Log.Debug("mqtt acl: invalid topic prefix",
			"client", cl.ID,
			"topic", topic)
		return false
	}

	topicDeviceID := parts[1]
	if topicDeviceID != cl.ID {
		h.Log.Warn("mqtt acl: cross-device access denied",
			"client", cl.ID,
			"topic_device", topicDeviceID,
			"topic", topic)
		return false
	}

	// Grace period restriction: expired certs can only use renew and cert topics
	if h.isInGracePeriod(cl) {
		subtopic := ""
		if len(parts) == 3 {
			subtopic = parts[2]
		}
		if subtopic != "renew" && subtopic != "cert" {
			h.Log.Warn("mqtt acl: grace period restricts to renew/cert only",
				"client", cl.ID,
				"topic", topic)
			return false
		}
	}

	return true
}
