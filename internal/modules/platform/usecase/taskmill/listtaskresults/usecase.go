package listtaskresults

import (
	"context"
	"go-enterprise-blueprint/internal/modules/platform/domain"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/taskmill/console"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	QueueName       *string    `query:"queue_name"`
	TaskGroupID     *string    `query:"task_group_id"`
	CompletedAfter  *time.Time `query:"completed_after"`
	CompletedBefore *time.Time `query:"completed_before"`
	Limit           int        `query:"limit"`
	Offset          int        `query:"offset"`
}

type Response struct {
	Content []console.TaskResult `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "list-task-results" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// List completed task results with filters
	results, err := uc.domainContainer.Console().ListResults(ctx, console.ListResultsParams{
		QueueName:       in.QueueName,
		TaskGroupID:     in.TaskGroupID,
		CompletedAfter:  in.CompletedAfter,
		CompletedBefore: in.CompletedBefore,
		Limit:           in.Limit,
		Offset:          in.Offset,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return task results
	return &Response{Content: results}, nil
}
