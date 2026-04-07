package getmysessions

import (
	"context"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct{}

type Response struct {
	Content []session.Session `json:"content"`
}

// UseCase implements "get-my-sessions" user action.
type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-my-sessions" }

func (uc *usecase) Execute(ctx context.Context, _ *Request) (*Response, error) {
	// Get user ID from authenticated user context
	userCtx := auth.MustUserContext(ctx)

	// List all sessions for the user
	sessions, err := uc.domainContainer.SessionRepo().List(ctx, session.Filter{
		UserID: &userCtx.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{Content: sessions}, nil
}
