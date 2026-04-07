package cleanuperrors

import (
	"context"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/audit"
	portalplatform "github.com/jaxongir1006/hire-ready-api/internal/portal/platform"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	CreatedBefore time.Time `json:"created_before" validate:"required"`
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

func (uc *usecase) OperationID() string { return "cleanup-errors" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Delete errors older than the specified timestamp
	count, err := uc.domainContainer.AlertErrorRepo().DeleteOlderThan(ctx, in.CreatedBefore)
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
