package org

import (
	"context"

	identityrepo "rootstock/web-server/repo/identity"
)

// Ops holds organization operations. Each method is one op.
type Ops struct {
	repo identityrepo.Repository
}

// NewOps creates org ops backed by the given identity repository.
func NewOps(repo identityrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// CreateOrg creates a new organization.
func (o *Ops) CreateOrg(ctx context.Context, input CreateOrgInput) (*Org, error) {
	result, err := o.repo.CreateOrg(ctx, toRepoCreateOrgInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoOrg(result), nil
}

// NestOrg creates a child organization under a parent.
func (o *Ops) NestOrg(ctx context.Context, input NestOrgInput) (*Org, error) {
	result, err := o.repo.NestOrg(ctx, toRepoNestOrgInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoOrg(result), nil
}

// DefineRole defines a role on a project.
func (o *Ops) DefineRole(ctx context.Context, input DefineRoleInput) (*Role, error) {
	result, err := o.repo.DefineRole(ctx, toRepoDefineRoleInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoRole(result), nil
}

// AssignRole grants roles to a user on a project.
func (o *Ops) AssignRole(ctx context.Context, input AssignRoleInput) (*UserGrant, error) {
	result, err := o.repo.AssignRole(ctx, toRepoAssignRoleInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoUserGrant(result), nil
}

// InviteUser creates a user and sends an invitation.
func (o *Ops) InviteUser(ctx context.Context, input InviteUserInput) (*InviteResult, error) {
	result, err := o.repo.InviteUser(ctx, toRepoInviteUserInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoInviteResult(result), nil
}

// CreateIdpUser creates a human user in the identity provider with a password.
func (o *Ops) CreateIdpUser(ctx context.Context, input CreateIdpUserInput) (*CreatedIdpUser, error) {
	result, err := o.repo.CreateUser(ctx, identityrepo.CreateHumanUserInput{
		Email:      input.Email,
		Password:   input.Password,
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
	})
	if err != nil {
		return nil, err
	}
	return &CreatedIdpUser{UserID: result.UserID, EmailCode: result.EmailCode}, nil
}

// VerifyEmail verifies a user's email address using the verification code.
func (o *Ops) VerifyEmail(ctx context.Context, input VerifyEmailInput) error {
	return o.repo.VerifyEmail(ctx, identityrepo.VerifyEmailInput{
		UserID:           input.UserID,
		VerificationCode: input.VerificationCode,
	})
}

func toRepoCreateOrgInput(in CreateOrgInput) identityrepo.CreateOrgInput {
	return identityrepo.CreateOrgInput{Name: in.Name}
}

func toRepoNestOrgInput(in NestOrgInput) identityrepo.NestOrgInput {
	return identityrepo.NestOrgInput{Name: in.Name, ParentOrgID: in.ParentOrgID}
}

func toRepoDefineRoleInput(in DefineRoleInput) identityrepo.DefineRoleInput {
	return identityrepo.DefineRoleInput{
		ProjectID:   in.ProjectID,
		RoleKey:     in.RoleKey,
		DisplayName: in.DisplayName,
	}
}

func toRepoAssignRoleInput(in AssignRoleInput) identityrepo.AssignRoleInput {
	return identityrepo.AssignRoleInput{
		UserID:    in.UserID,
		ProjectID: in.ProjectID,
		RoleKeys:  in.RoleKeys,
	}
}

func toRepoInviteUserInput(in InviteUserInput) identityrepo.InviteUserInput {
	return identityrepo.InviteUserInput{
		OrgID:      in.OrgID,
		Email:      in.Email,
		GivenName:  in.GivenName,
		FamilyName: in.FamilyName,
	}
}

func fromRepoOrg(r *identityrepo.Org) *Org {
	return &Org{ID: r.ID, Name: r.Name}
}

func fromRepoRole(r *identityrepo.Role) *Role {
	return &Role{ProjectID: r.ProjectID, RoleKey: r.RoleKey, DisplayName: r.DisplayName}
}

func fromRepoUserGrant(r *identityrepo.UserGrant) *UserGrant {
	return &UserGrant{
		UserGrantID: r.UserGrantID,
		UserID:      r.UserID,
		ProjectID:   r.ProjectID,
		RoleKeys:    r.RoleKeys,
	}
}

func fromRepoInviteResult(r *identityrepo.InviteResult) *InviteResult {
	return &InviteResult{UserID: r.UserID, EmailCode: r.EmailCode}
}
