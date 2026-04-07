package geterrorstats

import (
	"context"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/domain/alerterror"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	CreatedFrom *time.Time `query:"created_from"`
	CreatedTo   *time.Time `query:"created_to"`
}

type UseCase = ucdef.UserAction[*Request, *alerterror.Stats]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-error-stats" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*alerterror.Stats, error) {
	// Get error statistics for the given time range
	stats, err := uc.domainContainer.AlertErrorRepo().GetStats(ctx, alerterror.StatsFilter{
		CreatedFrom: in.CreatedFrom,
		CreatedTo:   in.CreatedTo,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return aggregated stats
	return stats, nil
}
