package triggerschedule

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	portalplatform "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/platform"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/taskmill/console"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	OperationID string `json:"operation_id" validate:"required"`
}

type UseCase = ucdef.UserAction[*Request, *struct{}]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "trigger-schedule" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*struct{}, error) {
	// Verify schedule exists by operation ID and enqueue task for immediate execution
	err := uc.domainContainer.Console().TriggerSchedule(ctx, in.OperationID)
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, console.CodeScheduleNotFound)
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
