package geterror

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain/alerterror"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	ID string `query:"id" validate:"required,uuid"`
}

type UseCase = ucdef.UserAction[*Request, *alerterror.Error]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-error" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*alerterror.Error, error) {
	// Find error by ID
	result, err := uc.domainContainer.AlertErrorRepo().Get(ctx, in.ID)
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, alerterror.CodeErrorNotFound)
	}

	// Return error details
	return result, nil
}
