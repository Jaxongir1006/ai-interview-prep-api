package purgequeue

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	portalplatform "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/platform"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	QueueName string `json:"queue_name" validate:"required"`
}

type UseCase = ucdef.UserAction[*Request, *struct{}]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "purge-queue" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*struct{}, error) {
	// Remove all non-DLQ tasks from the specified queue
	err := uc.domainContainer.Console().Purge(ctx, in.QueueName)
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

	return &struct{}{}, nil
}
