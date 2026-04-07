package setuserroles

import (
	"context"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/audit"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	UserID  string  `json:"user_id"  validate:"required"`
	RoleIDs []int64 `json:"role_ids" validate:"required"`
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

func (uc *usecase) OperationID() string { return "set-user-roles" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find user by ID
	_, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.UserID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// Validate all role IDs exist
	if len(in.RoleIDs) > 0 {
		roles, rolesErr := uc.domainContainer.RoleRepo().List(ctx, rbac.RoleFilter{
			IDs: in.RoleIDs,
		})
		if rolesErr != nil {
			return nil, errx.Wrap(rolesErr)
		}
		if len(roles) != len(in.RoleIDs) {
			return nil, errx.New(
				"one or more role IDs do not exist",
				errx.WithType(errx.T_NotFound),
				errx.WithCode(rbac.CodeRoleNotFound),
			)
		}
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Delete all existing role assignments for the user
	existing, err := uow.UserRole().List(ctx, rbac.UserRoleFilter{
		UserID: &in.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if len(existing) > 0 {
		err = uow.UserRole().BulkDelete(ctx, existing)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	// Create new role assignment records
	if len(in.RoleIDs) > 0 {
		newRoles := make([]rbac.UserRole, len(in.RoleIDs))
		for i, roleID := range in.RoleIDs {
			newRoles[i] = rbac.UserRole{
				UserID: in.UserID,
				RoleID: roleID,
			}
		}
		err = uow.UserRole().BulkCreate(ctx, newRoles)
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
