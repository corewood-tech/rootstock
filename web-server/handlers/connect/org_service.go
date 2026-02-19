package connect

import (
	"context"

	"connectrpc.com/connect"

	orgflows "rootstock/web-server/flows/org"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// OrgServiceHandler implements the OrgService Connect RPC interface.
type OrgServiceHandler struct {
	createOrg          *orgflows.CreateOrgFlow
	nestOrg            *orgflows.NestOrgFlow
	defineRole         *orgflows.DefineRoleFlow
	assignRole         *orgflows.AssignRoleFlow
	inviteUser         *orgflows.InviteUserFlow
}

// NewOrgServiceHandler creates the handler with all required flows.
func NewOrgServiceHandler(
	createOrg *orgflows.CreateOrgFlow,
	nestOrg *orgflows.NestOrgFlow,
	defineRole *orgflows.DefineRoleFlow,
	assignRole *orgflows.AssignRoleFlow,
	inviteUser *orgflows.InviteUserFlow,
) *OrgServiceHandler {
	return &OrgServiceHandler{
		createOrg:  createOrg,
		nestOrg:    nestOrg,
		defineRole: defineRole,
		assignRole: assignRole,
		inviteUser: inviteUser,
	}
}

func (h *OrgServiceHandler) CreateOrg(
	ctx context.Context,
	req *connect.Request[rootstockv1.CreateOrgRequest],
) (*connect.Response[rootstockv1.CreateOrgResponse], error) {
	result, err := h.createOrg.Run(ctx, orgflows.CreateOrgInput{
		Name: req.Msg.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.CreateOrgResponse{
		OrgId: result.ID,
		Name:  result.Name,
	}), nil
}

func (h *OrgServiceHandler) NestOrg(
	ctx context.Context,
	req *connect.Request[rootstockv1.NestOrgRequest],
) (*connect.Response[rootstockv1.NestOrgResponse], error) {
	result, err := h.nestOrg.Run(ctx, orgflows.NestOrgInput{
		Name:        req.Msg.GetName(),
		ParentOrgID: req.Msg.GetParentOrgId(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.NestOrgResponse{
		OrgId: result.ID,
		Name:  result.Name,
	}), nil
}

func (h *OrgServiceHandler) DefineRole(
	ctx context.Context,
	req *connect.Request[rootstockv1.DefineRoleRequest],
) (*connect.Response[rootstockv1.DefineRoleResponse], error) {
	result, err := h.defineRole.Run(ctx, orgflows.DefineRoleInput{
		ProjectID:   req.Msg.GetProjectId(),
		RoleKey:     req.Msg.GetRoleKey(),
		DisplayName: req.Msg.GetDisplayName(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.DefineRoleResponse{
		ProjectId:   result.ProjectID,
		RoleKey:     result.RoleKey,
		DisplayName: result.DisplayName,
	}), nil
}

func (h *OrgServiceHandler) AssignRole(
	ctx context.Context,
	req *connect.Request[rootstockv1.AssignRoleRequest],
) (*connect.Response[rootstockv1.AssignRoleResponse], error) {
	result, err := h.assignRole.Run(ctx, orgflows.AssignRoleInput{
		UserID:    req.Msg.GetUserId(),
		ProjectID: req.Msg.GetProjectId(),
		RoleKeys:  req.Msg.GetRoleKeys(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.AssignRoleResponse{
		UserGrantId: result.UserGrantID,
		UserId:      result.UserID,
		ProjectId:   result.ProjectID,
		RoleKeys:    result.RoleKeys,
	}), nil
}

func (h *OrgServiceHandler) InviteUser(
	ctx context.Context,
	req *connect.Request[rootstockv1.InviteUserRequest],
) (*connect.Response[rootstockv1.InviteUserResponse], error) {
	result, err := h.inviteUser.Run(ctx, orgflows.InviteUserInput{
		OrgID:      req.Msg.GetOrgId(),
		Email:      req.Msg.GetEmail(),
		GivenName:  req.Msg.GetGivenName(),
		FamilyName: req.Msg.GetFamilyName(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&rootstockv1.InviteUserResponse{
		UserId:    result.UserID,
		EmailCode: result.EmailCode,
	}), nil
}
