package disableuser

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	ID string `json:"id" validate:"required"`
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

func (uc *usecase) OperationID() string { return "disable-user" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find user by ID
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.ID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// Check if user is not already disabled
	if !u.IsActive {
		return nil, errx.New(
			"user is already disabled",
			errx.WithType(errx.T_Validation),
			errx.WithCode(user.CodeUserAlreadyDisabled),
		)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Set user is_active to false
	u.IsActive = false
	_, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Delete all user sessions (force logout)
	sessions, err := uow.Session().List(ctx, session.Filter{UserID: &u.ID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if len(sessions) > 0 {
		err = uow.Session().BulkDelete(ctx, sessions)
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
