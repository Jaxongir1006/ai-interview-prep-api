package cleanupresults

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	portalplatform "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/platform"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/taskmill/console"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	CompletedBefore time.Time `json:"completed_before" validate:"required"`
	QueueName       *string   `json:"queue_name"`
}

type Response struct {
	DeletedCount int64 `json:"deleted_count"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "cleanup-results" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Delete old task results matching criteria
	count, err := uc.domainContainer.Console().CleanupResults(ctx, console.CleanupResultsParams{
		CompletedBefore: in.CompletedBefore,
		QueueName:       in.QueueName,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Record audit log
	err = uc.portalContainer.Audit().Log(uow.Lend(), audit.Action{
		Module: portalplatform.ModuleName, OperationID: uc.OperationID(), Payload: in,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return count of deleted records
	return &Response{DeletedCount: count}, nil
}
