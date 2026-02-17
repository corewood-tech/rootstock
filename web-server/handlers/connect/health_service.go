package connect

import (
	"context"

	"connectrpc.com/connect"

	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

type HealthServiceHandler struct{}

func NewHealthServiceHandler() *HealthServiceHandler {
	return &HealthServiceHandler{}
}

func (h *HealthServiceHandler) Check(
	ctx context.Context,
	req *connect.Request[rootstockv1.CheckRequest],
) (*connect.Response[rootstockv1.CheckResponse], error) {
	return connect.NewResponse(&rootstockv1.CheckResponse{
		Status: "ok",
	}), nil
}
