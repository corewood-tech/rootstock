package connect

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	userflows "rootstock/web-server/flows/user"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
	"rootstock/web-server/auth"
)

// UserServiceHandler implements the UserService Connect RPC interface.
type UserServiceHandler struct {
	registerUser *userflows.RegisterUserFlow
	getUser      *userflows.GetUserFlow
}

// NewUserServiceHandler creates the handler with all required flows.
func NewUserServiceHandler(
	registerUser *userflows.RegisterUserFlow,
	getUser *userflows.GetUserFlow,
) *UserServiceHandler {
	return &UserServiceHandler{
		registerUser: registerUser,
		getUser:      getUser,
	}
}

func (h *UserServiceHandler) RegisterUser(
	ctx context.Context,
	req *connect.Request[rootstockv1.RegisterUserRequest],
) (*connect.Response[rootstockv1.RegisterUserResponse], error) {
	idpID, ok := auth.SubjectFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no authenticated subject"))
	}

	result, err := h.registerUser.Run(ctx, userflows.RegisterUserInput{
		IdpID:    idpID,
		UserType: req.Msg.GetUserType(),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.RegisterUserResponse{
		User: userToProto(result),
	}), nil
}

func (h *UserServiceHandler) GetMe(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetMeRequest],
) (*connect.Response[rootstockv1.GetMeResponse], error) {
	idpID, ok := auth.SubjectFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no authenticated subject"))
	}

	result, err := h.getUser.Run(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	return connect.NewResponse(&rootstockv1.GetMeResponse{
		User: userToProto(result),
	}), nil
}

func userToProto(u *userflows.User) *rootstockv1.UserProto {
	return &rootstockv1.UserProto{
		Id:        u.ID,
		UserType:  u.UserType,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
