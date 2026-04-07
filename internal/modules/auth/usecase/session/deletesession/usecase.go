package deletesession

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	SessionID int64 `json:"session_id" validate:"required"`
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

func (uc *usecase) OperationID() string { return "delete-session" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find session by ID
	s, err := uc.domainContainer.SessionRepo().Get(ctx, session.Filter{
		ID: &in.SessionID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, session.CodeSessionNotFound)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Delete session
	err = uow.Session().Delete(ctx, s)
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
