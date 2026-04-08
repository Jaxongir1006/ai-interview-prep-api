package listschedules

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/taskmill/console"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	QueueName *string `query:"queue_name"`
}

type Response struct {
	Content []console.ScheduleInfo `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "list-schedules" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// List cron schedules, optionally filtered by queue
	schedules, err := uc.domainContainer.Console().ListSchedules(ctx, in.QueueName)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return schedules
	return &Response{Content: schedules}, nil
}
