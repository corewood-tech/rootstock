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
	registerUser       *userflows.RegisterUserFlow
	getUser            *userflows.GetUserFlow
	login              *userflows.LoginFlow
	logout             *userflows.LogoutFlow
	registerResearcher *userflows.RegisterResearcherFlow
}

// NewUserServiceHandler creates the handler with all required flows.
func NewUserServiceHandler(
	registerUser *userflows.RegisterUserFlow,
	getUser *userflows.GetUserFlow,
	login *userflows.LoginFlow,
	logout *userflows.LogoutFlow,
	registerResearcher *userflows.RegisterResearcherFlow,
) *UserServiceHandler {
	return &UserServiceHandler{
		registerUser:       registerUser,
		getUser:            getUser,
		login:              login,
		logout:             logout,
		registerResearcher: registerResearcher,
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

func (h *UserServiceHandler) Login(
	ctx context.Context,
	req *connect.Request[rootstockv1.LoginRequest],
) (*connect.Response[rootstockv1.LoginResponse], error) {
	result, err := h.login.Run(ctx, userflows.LoginInput{
		Email:    req.Msg.GetEmail(),
		Password: req.Msg.GetPassword(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("login failed: %w", err))
	}

	return connect.NewResponse(&rootstockv1.LoginResponse{
		SessionId:    result.SessionID,
		SessionToken: result.SessionToken,
		User:         flowUserToProto(&result.User),
	}), nil
}

func (h *UserServiceHandler) Logout(
	ctx context.Context,
	req *connect.Request[rootstockv1.LogoutRequest],
) (*connect.Response[rootstockv1.LogoutResponse], error) {
	sessionID, ok := auth.SessionIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no session"))
	}
	sessionToken, ok := auth.SessionTokenFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no session token"))
	}

	err := h.logout.Run(ctx, userflows.LogoutInput{
		SessionID:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("logout failed: %w", err))
	}

	return connect.NewResponse(&rootstockv1.LogoutResponse{}), nil
}

func (h *UserServiceHandler) RegisterResearcher(
	ctx context.Context,
	req *connect.Request[rootstockv1.RegisterResearcherRequest],
) (*connect.Response[rootstockv1.RegisterResearcherResponse], error) {
	result, err := h.registerResearcher.Run(ctx, userflows.RegisterResearcherInput{
		Email:      req.Msg.GetEmail(),
		Password:   req.Msg.GetPassword(),
		GivenName:  req.Msg.GetGivenName(),
		FamilyName: req.Msg.GetFamilyName(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("registration failed: %w", err))
	}

	return connect.NewResponse(&rootstockv1.RegisterResearcherResponse{
		SessionId:    result.SessionID,
		SessionToken: result.SessionToken,
		User:         flowUserToProto(&result.User),
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

func flowUserToProto(u *userflows.User) *rootstockv1.UserProto {
	return &rootstockv1.UserProto{
		Id:        u.ID,
		UserType:  u.UserType,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
