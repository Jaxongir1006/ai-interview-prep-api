package updaterole

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
	ID   int64  `json:"id"   validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type UseCase = ucdef.UserAction[*Request, *rbac.Role]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "update-role" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*rbac.Role, error) {
	// Find role by ID
	role, err := uc.domainContainer.RoleRepo().Get(ctx, rbac.RoleFilter{
		ID: &in.ID,
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

	// Update role name
	role.Name = in.Name
	role, err = uow.Role().Update(ctx, role)
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, rbac.CodeRoleNameConflict)
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

	return role, nil
}
