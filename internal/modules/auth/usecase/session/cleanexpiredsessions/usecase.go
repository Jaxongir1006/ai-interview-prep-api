package cleanexpiredsessions

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/ucdef"
)

// Payload is empty for this scheduled task.
// The task operates on all expired sessions without external parameters.
type Payload struct{}

type UseCase = ucdef.AsyncTask[*Payload]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{
		domainContainer,
	}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "clean-expired-sessions" }

func (uc *usecase) Execute(ctx context.Context, _ *Payload) error {
	count, err := uc.domainContainer.SessionRepo().DeleteExpired(ctx)
	if err != nil {
		return errx.Wrap(err)
	}

	logger.
		WithContext(ctx).
		Named("auth_usecase").
		With("operation_id", uc.OperationID()).
		With("deleted_count", count).
		Info("expired sessions cleaned up")

	return nil
}
