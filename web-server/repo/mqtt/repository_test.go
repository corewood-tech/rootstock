package mqtt

import (
	"context"
	"sync"
	"testing"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/packets"
)

func setupBroker(t *testing.T) *mochi.Server {
	t.Helper()

	server := mochi.New(&mochi.Options{InlineClient: true})
	// AllowHook permits all connections (no mTLS in unit test)
	if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
		t.Fatalf("add allow hook: %v", err)
	}
	go server.Serve()
	t.Cleanup(func() { server.Close() })

	// Wait for broker to start
	time.Sleep(100 * time.Millisecond)
	return server
}

func TestPushDeviceConfig(t *testing.T) {
	server := setupBroker(t)
	repo := NewRepository(server)
	defer repo.Shutdown()

	var mu sync.Mutex
	var received []byte

	// Subscribe to the expected config topic
	err := server.Subscribe("rootstock/device-001/config", 1, func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
		mu.Lock()
		received = append([]byte{}, pk.Payload...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	ctx := context.Background()
	payload := []byte(`{"campaign_id":"camp-1","interval":60}`)

	if err := repo.PushDeviceConfig(ctx, PushConfigInput{
		DeviceID: "device-001",
		Payload:  payload,
	}); err != nil {
		t.Fatalf("PushDeviceConfig(): %v", err)
	}

	// Wait for delivery
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if string(received) != string(payload) {
		t.Errorf("received %q, want %q", received, payload)
	}
}

func TestPublishToDevice(t *testing.T) {
	server := setupBroker(t)
	repo := NewRepository(server)
	defer repo.Shutdown()

	var mu sync.Mutex
	var received []byte

	err := server.Subscribe("rootstock/device-002/cert", 1, func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
		mu.Lock()
		received = append([]byte{}, pk.Payload...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	ctx := context.Background()
	certPEM := []byte("-----BEGIN CERTIFICATE-----\nfake\n-----END CERTIFICATE-----")

	if err := repo.PublishToDevice(ctx, PublishInput{
		Topic:   "rootstock/device-002/cert",
		Payload: certPEM,
		QoS:     1,
	}); err != nil {
		t.Fatalf("PublishToDevice(): %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if string(received) != string(certPEM) {
		t.Errorf("received %q, want %q", received, certPEM)
	}
}

func TestShutdown(t *testing.T) {
	server := setupBroker(t)
	repo := NewRepository(server)

	// Should not hang
	done := make(chan struct{})
	go func() {
		repo.Shutdown()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("Shutdown() timed out")
	}
}
