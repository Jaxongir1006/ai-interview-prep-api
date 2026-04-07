package setrolepermissions

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	RoleID      int64    `json:"role_id"     validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

type Response struct{}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "set-role-permissions" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find role by ID
	_, err := uc.domainContainer.RoleRepo().Get(ctx, rbac.RoleFilter{
		ID: &in.RoleID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, rbac.CodeRoleNotFound)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Delete all existing permissions for the role
	existing, err := uow.RolePermission().List(ctx, rbac.RolePermissionFilter{
		RoleID: &in.RoleID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if len(existing) > 0 {
		err = uow.RolePermission().BulkDelete(ctx, existing)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	// Create new permission records
	if len(in.Permissions) > 0 {
		newPerms := make([]rbac.RolePermission, len(in.Permissions))
		for i, p := range in.Permissions {
			newPerms[i] = rbac.RolePermission{
				RoleID:     in.RoleID,
				Permission: p,
			}
		}
		err = uow.RolePermission().BulkCreate(ctx, newPerms)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	// Record audit log
	err = uc.portalContainer.Audit().Log(uow.Lend(), audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: in,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{}, nil
}
