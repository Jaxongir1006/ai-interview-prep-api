package purgedlq

import (
	"context"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/audit"
	portalplatform "github.com/jaxongir1006/hire-ready-api/internal/portal/platform"

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

func (uc *usecase) OperationID() string { return "purge-dlq" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*struct{}, error) {
	// Remove all DLQ tasks from the specified queue
	err := uc.domainContainer.Console().PurgeDLQ(ctx, in.QueueName)
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
