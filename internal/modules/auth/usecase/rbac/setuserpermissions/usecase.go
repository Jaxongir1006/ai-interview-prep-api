package setuserpermissions

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/rbac"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	UserID      string   `json:"user_id"     validate:"required"`
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

func (uc *usecase) OperationID() string { return "set-user-permissions" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find user by ID
	_, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.UserID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Delete all existing direct permissions for the user
	existing, err := uow.UserPermission().List(ctx, rbac.UserPermissionFilter{
		UserID: &in.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if len(existing) > 0 {
		err = uow.UserPermission().BulkDelete(ctx, existing)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	// Create new permission records
	if len(in.Permissions) > 0 {
		newPerms := make([]rbac.UserPermission, len(in.Permissions))
		for i, p := range in.Permissions {
			newPerms[i] = rbac.UserPermission{
				UserID:     in.UserID,
				Permission: p,
			}
		}
		err = uow.UserPermission().BulkCreate(ctx, newPerms)
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
