package connect

import (
	"context"
	"time"

	"connectrpc.com/connect"

	deviceflows "rootstock/web-server/flows/device"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// DeviceServiceHandler implements the DeviceService Connect RPC interface.
type DeviceServiceHandler struct {
	getDevice        *deviceflows.GetDeviceFlow
	revokeDevice     *deviceflows.RevokeDeviceFlow
	reinstateDevice  *deviceflows.ReinstateDeviceFlow
}

// NewDeviceServiceHandler creates the handler with all required flows.
func NewDeviceServiceHandler(
	getDevice *deviceflows.GetDeviceFlow,
	revokeDevice *deviceflows.RevokeDeviceFlow,
	reinstateDevice *deviceflows.ReinstateDeviceFlow,
) *DeviceServiceHandler {
	return &DeviceServiceHandler{
		getDevice:       getDevice,
		revokeDevice:    revokeDevice,
		reinstateDevice: reinstateDevice,
	}
}

func (h *DeviceServiceHandler) GetDevice(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetDeviceRequest],
) (*connect.Response[rootstockv1.GetDeviceResponse], error) {
	device, err := h.getDevice.Run(ctx, deviceflows.GetDeviceInput{
		DeviceID: req.Msg.GetDeviceId(),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.GetDeviceResponse{
		Device: deviceToProto(device),
	}), nil
}

func (h *DeviceServiceHandler) RevokeDevice(
	ctx context.Context,
	req *connect.Request[rootstockv1.RevokeDeviceRequest],
) (*connect.Response[rootstockv1.RevokeDeviceResponse], error) {
	if err := h.revokeDevice.Run(ctx, deviceflows.RevokeDeviceInput{
		DeviceID: req.Msg.GetDeviceId(),
	}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.RevokeDeviceResponse{}), nil
}

func (h *DeviceServiceHandler) ReinstateDevice(
	ctx context.Context,
	req *connect.Request[rootstockv1.ReinstateDeviceRequest],
) (*connect.Response[rootstockv1.ReinstateDeviceResponse], error) {
	if err := h.reinstateDevice.Run(ctx, deviceflows.ReinstateDeviceInput{
		DeviceID: req.Msg.GetDeviceId(),
	}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.ReinstateDeviceResponse{}), nil
}

func deviceToProto(d *deviceflows.Device) *rootstockv1.DeviceProto {
	proto := &rootstockv1.DeviceProto{
		Id:              d.ID,
		OwnerId:         d.OwnerID,
		Status:          d.Status,
		Class:           d.Class,
		FirmwareVersion: d.FirmwareVersion,
		Tier:            int32(d.Tier),
		Sensors:         d.Sensors,
		CertSerial:      d.CertSerial,
		CreatedAt:       d.CreatedAt.Format(time.RFC3339),
	}
	return proto
}
