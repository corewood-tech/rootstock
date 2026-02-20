package connect

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	securityflows "rootstock/web-server/flows/security"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// AdminServiceHandler implements the AdminService Connect RPC interface.
type AdminServiceHandler struct {
	securityResponse *securityflows.SecurityResponseFlow
}

// NewAdminServiceHandler creates the handler with all required flows.
func NewAdminServiceHandler(securityResponse *securityflows.SecurityResponseFlow) *AdminServiceHandler {
	return &AdminServiceHandler{
		securityResponse: securityResponse,
	}
}

func (h *AdminServiceHandler) SuspendByClass(
	ctx context.Context,
	req *connect.Request[rootstockv1.SuspendByClassRequest],
) (*connect.Response[rootstockv1.SuspendByClassResponse], error) {
	windowStart, err := time.Parse(time.RFC3339, req.Msg.GetWindowStart())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse window_start: %w", err))
	}
	windowEnd, err := time.Parse(time.RFC3339, req.Msg.GetWindowEnd())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse window_end: %w", err))
	}

	result, err := h.securityResponse.Run(ctx, securityflows.SecurityResponseInput{
		Class:       req.Msg.GetDeviceClass(),
		FirmwareMin: req.Msg.GetFirmwareMin(),
		FirmwareMax: req.Msg.GetFirmwareMax(),
		WindowStart: windowStart,
		WindowEnd:   windowEnd,
		Reason:      req.Msg.GetReason(),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.SuspendByClassResponse{
		SuspendedCount:      int32(result.SuspendedCount),
		QuarantinedReadings: result.QuarantinedReadings,
		NotifiedScitizens:   int32(result.NotifiedScitizens),
	}), nil
}
