package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	deviceflows "rootstock/web-server/flows/device"
	readingflows "rootstock/web-server/flows/reading"
	scoreflows "rootstock/web-server/flows/score"
	"rootstock/web-server/global/observability"
	mqttrepo "rootstock/web-server/repo/mqtt"
)

// MQTTFlows holds the flows that MQTT inline subscriptions invoke.
type MQTTFlows struct {
	IngestReading        *readingflows.IngestReadingFlow
	RenewCert            *deviceflows.RenewCertFlow
	RefreshScitizenScore *scoreflows.RefreshScitizenScoreFlow
}

// ReadingPayload is the JSON payload published by devices on telemetry topics.
type ReadingPayload struct {
	Value           float64   `json:"value"`
	Timestamp       time.Time `json:"timestamp"`
	Geolocation     string    `json:"geolocation,omitempty"`
	FirmwareVersion string    `json:"firmware_version"`
	CertSerial      string    `json:"cert_serial"`
}

// SetupMQTTSubscriptions registers inline subscriptions on the embedded broker
// that route MQTT messages to the appropriate flows. Call after all flows are
// constructed but before server.Serve().
func SetupMQTTSubscriptions(ctx context.Context, server *mochi.Server, flows *MQTTFlows) error {
	logger := observability.GetLogger("mqtt-subscriptions")

	// Telemetry: rootstock/+/data/+
	telemetryTopic := fmt.Sprintf("%s/+/data/+", mqttrepo.TopicPrefix)
	if err := server.Subscribe(telemetryTopic, 1, func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
		segments := strings.Split(pk.TopicName, "/")
		if len(segments) < 4 {
			logger.Error(ctx, "telemetry: unexpected topic format", map[string]interface{}{
				"topic": pk.TopicName,
			})
			return
		}
		deviceID := segments[1]
		campaignID := segments[3]

		var payload ReadingPayload
		if err := json.Unmarshal(pk.Payload, &payload); err != nil {
			logger.Error(ctx, "telemetry: invalid payload JSON", map[string]interface{}{
				"device_id": deviceID,
				"error":     err.Error(),
			})
			return
		}

		input := readingflows.IngestReadingInput{
			DeviceID:        deviceID,
			CampaignID:      campaignID,
			Value:           payload.Value,
			Timestamp:       payload.Timestamp,
			Geolocation:     payload.Geolocation,
			FirmwareVersion: payload.FirmwareVersion,
			CertSerial:      payload.CertSerial,
		}

		result, err := flows.IngestReading.Run(ctx, input)
		if err != nil {
			logger.Error(ctx, "telemetry: ingest reading failed", map[string]interface{}{
				"device_id":   deviceID,
				"campaign_id": campaignID,
				"error":       err.Error(),
			})
			return
		}

		if result.Status == "accepted" {
			if _, err := flows.RefreshScitizenScore.Run(ctx, scoreflows.RefreshScitizenScoreInput{
				DeviceID: deviceID,
			}); err != nil {
				logger.Error(ctx, "telemetry: refresh scitizen score failed", map[string]interface{}{
					"device_id": deviceID,
					"error":     err.Error(),
				})
			}
		}
	}); err != nil {
		return fmt.Errorf("subscribe telemetry: %w", err)
	}

	// Renewal: rootstock/+/renew
	renewTopic := fmt.Sprintf("%s/+/renew", mqttrepo.TopicPrefix)
	if err := server.Subscribe(renewTopic, 1, func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
		segments := strings.Split(pk.TopicName, "/")
		if len(segments) < 3 {
			logger.Error(ctx, "renew: unexpected topic format", map[string]interface{}{
				"topic": pk.TopicName,
			})
			return
		}
		deviceID := segments[1]

		result, err := flows.RenewCert.Run(ctx, deviceflows.RenewCertInput{
			DeviceID: deviceID,
			CSR:      pk.Payload,
		})
		if err != nil {
			logger.Error(ctx, "renew: certificate renewal failed", map[string]interface{}{
				"device_id": deviceID,
				"error":     err.Error(),
			})
			return
		}

		certTopic := fmt.Sprintf("%s/%s/cert", mqttrepo.TopicPrefix, deviceID)
		if err := server.Publish(certTopic, result.CertPEM, false, 1); err != nil {
			logger.Error(ctx, "renew: publish cert response failed", map[string]interface{}{
				"device_id": deviceID,
				"error":     err.Error(),
			})
		}
	}); err != nil {
		return fmt.Errorf("subscribe renew: %w", err)
	}

	logger.Info(ctx, "mqtt subscriptions registered", map[string]interface{}{
		"telemetry": telemetryTopic,
		"renew":     renewTopic,
	})

	return nil
}
