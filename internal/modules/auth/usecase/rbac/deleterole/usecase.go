package deleterole

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
	ID int64 `json:"id" validate:"required"`
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

func (uc *usecase) OperationID() string { return "delete-role" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find role by ID
	role, err := uc.domainContainer.RoleRepo().Get(ctx, rbac.RoleFilter{
		ID: &in.ID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, rbac.CodeRoleNotFound)
	}

	// Check role has no assigned users
	hasUsers, err := uc.domainContainer.UserRoleRepo().Exists(ctx, rbac.UserRoleFilter{
		RoleID: &role.ID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if hasUsers {
		return nil, errx.New(
			"role has assigned users and cannot be deleted",
			errx.WithType(errx.T_Validation),
			errx.WithCode(rbac.CodeRoleHasAssignedUsers),
		)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Delete all role permissions
	perms, err := uow.RolePermission().List(ctx, rbac.RolePermissionFilter{
		RoleID: &role.ID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if len(perms) > 0 {
		err = uow.RolePermission().BulkDelete(ctx, perms)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	// Delete role
	err = uow.Role().Delete(ctx, role)
	if err != nil {
		return nil, errx.Wrap(err)
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
