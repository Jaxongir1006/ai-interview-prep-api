package deletemysession

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	SessionID int64 `json:"session_id" validate:"required"`
}

type Response struct{}

// UseCase implements "delete-my-session" user action.
type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "delete-my-session" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Get user ID from authenticated user context
	userCtx := auth.MustUserContext(ctx)

	// Find target session by ID
	target, err := uc.domainContainer.SessionRepo().Get(ctx, session.Filter{
		ID: &in.SessionID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, session.CodeSessionNotFound)
	}

	// Verify target session belongs to the current user
	if target.UserID != userCtx.UserID {
		return nil, errx.New(
			"session does not belong to the current user",
			errx.WithType(errx.T_NotFound),
			errx.WithCode(session.CodeSessionNotFound),
		)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Delete target session
	err = uow.Session().Delete(ctx, target)
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
