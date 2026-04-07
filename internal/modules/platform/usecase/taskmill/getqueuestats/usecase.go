package getqueuestats

import (
	"context"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/domain"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	QueueName string `query:"queue_name" validate:"required"`
}

type Response struct {
	QueueName   string     `json:"queue_name"`
	Total       int64      `json:"total"`
	Available   int64      `json:"available"`
	InFlight    int64      `json:"in_flight"`
	Scheduled   int64      `json:"scheduled"`
	InDLQ       int64      `json:"in_dlq"`
	OldestTask  *time.Time `json:"oldest_task"`
	AvgAttempts float64    `json:"avg_attempts"`
	P95Attempts float64    `json:"p95_attempts"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-queue-stats" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Get stats for the specified queue
	stats, err := uc.domainContainer.Console().Stats(ctx, in.QueueName)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return queue statistics
	return &Response{
		QueueName:   stats.QueueName,
		Total:       stats.Total,
		Available:   stats.Available,
		InFlight:    stats.InFlight,
		Scheduled:   stats.Scheduled,
		InDLQ:       stats.InDLQ,
		OldestTask:  stats.OldestTask,
		AvgAttempts: stats.AvgAttempts,
		P95Attempts: stats.P95Attempts,
	}, nil
}
