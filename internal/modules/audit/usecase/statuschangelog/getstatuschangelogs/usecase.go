package getstatuschangelogs

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/statuschangelog"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	From        string  `query:"from"          validate:"required"`
	To          string  `query:"to"            validate:"required"`
	EntityType  *string `query:"entity_type"`
	EntityID    *string `query:"entity_id"`
	ActionLogID *int64  `query:"action_log_id"`
	Cursor      *int64  `query:"cursor"`
	Limit       *int    `query:"limit"`
}

type UseCase = ucdef.UserAction[*Request, []statuschangelog.StatusChangeLog]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-status-change-logs" }

func (uc *usecase) Execute(
	ctx context.Context,
	in *Request,
) ([]statuschangelog.StatusChangeLog, error) {
	// Parse and validate time range parameters
	from, err := time.Parse(time.RFC3339, in.From)
	if err != nil {
		return nil, errx.New(
			"invalid 'from' time format, expected RFC3339",
			errx.WithType(errx.T_Validation),
			errx.WithCode("INVALID_TIME_FORMAT"),
			errx.WithDetails(errx.D{
				"from": from,
			}),
		)
	}

	to, err := time.Parse(time.RFC3339, in.To)
	if err != nil {
		return nil, errx.New(
			"invalid 'to' time format, expected RFC3339",
			errx.WithType(errx.T_Validation),
			errx.WithCode("INVALID_TIME_FORMAT"),
			errx.WithDetails(errx.D{
				"to": to,
			}),
		)
	}

	if !from.Before(to) {
		return nil, errx.New(
			"'from' must be before 'to'",
			errx.WithType(errx.T_Validation),
			errx.WithCode("INVALID_TIME_RANGE"),
			errx.WithDetails(errx.D{
				"from": from,
				"to":   to,
			}),
		)
	}

	// List status change logs matching filter criteria
	logs, err := uc.domainContainer.StatusChangeLogRepo().List(ctx, statuschangelog.Filter{
		EntityType:  in.EntityType,
		EntityID:    in.EntityID,
		ActionLogID: in.ActionLogID,
		CreatedFrom: &from,
		CreatedTo:   &to,
		Cursor:      in.Cursor,
		Limit:       in.Limit,
	})

	// Return status change logs
	return logs, errx.Wrap(err)
}
