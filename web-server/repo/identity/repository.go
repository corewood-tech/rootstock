package identity

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel-go/v3/pkg/client"
	mgmt "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	objectv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object/v2"
	orgv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"
	userv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"rootstock/web-server/config"
)

type zitadelRepository struct {
	client *client.Client
}

// NewRepository creates an identity repository backed by Zitadel.
func NewRepository(ctx context.Context, cfg config.ZitadelConfig) (Repository, error) {
	z := zitadel.New(
		cfg.Host,
		zitadel.WithInsecure(fmt.Sprintf("%d", cfg.Port)),
	)

	c, err := client.New(ctx, z, client.WithAuth(client.PAT(cfg.ServiceUserPAT)))
	if err != nil {
		return nil, fmt.Errorf("create zitadel client: %w", err)
	}

	return &zitadelRepository{client: c}, nil
}

func (r *zitadelRepository) CreateOrg(ctx context.Context, input CreateOrgInput) (*Org, error) {
	resp, err := r.client.OrganizationServiceV2().AddOrganization(ctx, &orgv2.AddOrganizationRequest{
		Name: input.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("create org: %w", err)
	}
	return &Org{
		ID:   resp.GetOrganizationId(),
		Name: input.Name,
	}, nil
}

func (r *zitadelRepository) NestOrg(ctx context.Context, input NestOrgInput) (*Org, error) {
	// Zitadel v2 doesn't have explicit nesting — we create the org and store
	// the parent relationship via org metadata.
	resp, err := r.client.OrganizationServiceV2().AddOrganization(ctx, &orgv2.AddOrganizationRequest{
		Name: input.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("create nested org: %w", err)
	}

	orgID := resp.GetOrganizationId()

	// Store parent org ID as metadata on the new org.
	_, err = r.client.OrganizationServiceV2().SetOrganizationMetadata(ctx, &orgv2.SetOrganizationMetadataRequest{
		OrganizationId: orgID,
		Metadata: []*orgv2.Metadata{
			{
				Key:   "parent_org_id",
				Value: []byte(input.ParentOrgID),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("set parent org metadata: %w", err)
	}

	return &Org{
		ID:   orgID,
		Name: input.Name,
	}, nil
}

func (r *zitadelRepository) DefineRole(ctx context.Context, input DefineRoleInput) (*Role, error) {
	_, err := r.client.ManagementService().AddProjectRole(ctx, &mgmt.AddProjectRoleRequest{
		ProjectId:   input.ProjectID,
		RoleKey:     input.RoleKey,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("define role: %w", err)
	}
	return &Role{
		ProjectID:   input.ProjectID,
		RoleKey:     input.RoleKey,
		DisplayName: input.DisplayName,
	}, nil
}

func (r *zitadelRepository) AssignRole(ctx context.Context, input AssignRoleInput) (*UserGrant, error) {
	resp, err := r.client.ManagementService().AddUserGrant(ctx, &mgmt.AddUserGrantRequest{
		UserId:    input.UserID,
		ProjectId: input.ProjectID,
		RoleKeys:  input.RoleKeys,
	})
	if err != nil {
		return nil, fmt.Errorf("assign role: %w", err)
	}
	return &UserGrant{
		UserGrantID: resp.GetUserGrantId(),
		UserID:      input.UserID,
		ProjectID:   input.ProjectID,
		RoleKeys:    input.RoleKeys,
	}, nil
}

func (r *zitadelRepository) InviteUser(ctx context.Context, input InviteUserInput) (*InviteResult, error) {
	resp, err := r.client.UserServiceV2().AddHumanUser(ctx, &userv2.AddHumanUserRequest{
		Organization: &objectv2.Organization{
			Org: &objectv2.Organization_OrgId{OrgId: input.OrgID},
		},
		Profile: &userv2.SetHumanProfile{
			GivenName:  input.GivenName,
			FamilyName: input.FamilyName,
		},
		Email: &userv2.SetHumanEmail{
			Email: input.Email,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("invite user: %w", err)
	}
	return &InviteResult{
		UserID:    resp.GetUserId(),
		EmailCode: resp.GetEmailCode(),
	}, nil
}

func (r *zitadelRepository) CreateUser(ctx context.Context, input CreateHumanUserInput) (*CreatedUser, error) {
	resp, err := r.client.UserServiceV2().AddHumanUser(ctx, &userv2.AddHumanUserRequest{
		Profile: &userv2.SetHumanProfile{
			GivenName:  input.GivenName,
			FamilyName: input.FamilyName,
		},
		Email: &userv2.SetHumanEmail{
			Email:        input.Email,
			Verification: &userv2.SetHumanEmail_IsVerified{IsVerified: true},
		},
		PasswordType: &userv2.AddHumanUserRequest_Password{
			Password: &userv2.Password{
				Password:       input.Password,
				ChangeRequired: false,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &CreatedUser{
		UserID: resp.GetUserId(),
	}, nil
}

func (r *zitadelRepository) Shutdown() {
	r.client.Close()
}
