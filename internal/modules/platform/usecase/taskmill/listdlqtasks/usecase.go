package listdlqtasks

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/taskmill/console"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	QueueName   *string    `query:"queue_name"`
	OperationID *string    `query:"operation_id"`
	DLQAfter    *time.Time `query:"dlq_after"`
	DLQBefore   *time.Time `query:"dlq_before"`
	Limit       int        `query:"limit"`
	Offset      int        `query:"offset"`
}

type Response struct {
	Content []console.DLQTask `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "list-dlq-tasks" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// List DLQ tasks with filters
	tasks, err := uc.domainContainer.Console().ListDLQTasks(ctx, console.ListDLQTasksParams{
		QueueName:   in.QueueName,
		OperationID: in.OperationID,
		DLQAfter:    in.DLQAfter,
		DLQBefore:   in.DLQBefore,
		Limit:       in.Limit,
		Offset:      in.Offset,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return DLQ tasks
	return &Response{Content: tasks}, nil
}
