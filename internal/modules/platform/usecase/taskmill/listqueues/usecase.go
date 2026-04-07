package listqueues

import (
	"context"
	"go-enterprise-blueprint/internal/modules/platform/domain"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct{}

type Response struct {
	Content []string `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "list-queues" }

func (uc *usecase) Execute(ctx context.Context, _ *Request) (*Response, error) {
	// List all distinct queue names
	queues, err := uc.domainContainer.Console().ListQueues(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return queue names
	return &Response{Content: queues}, nil
}
