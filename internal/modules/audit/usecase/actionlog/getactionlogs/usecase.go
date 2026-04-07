package getactionlogs

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/actionlog"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	From        string   `query:"from"         validate:"required"`
	To          string   `query:"to"           validate:"required"`
	Module      *string  `query:"module"`
	OperationID *string  `query:"operation_id"`
	UserID      *string  `query:"user_id"`
	Tags        []string `query:"tags"`
	GroupKey    *string  `query:"group_key"`
	Cursor      *int64   `query:"cursor"`
	Limit       *int     `query:"limit"`
}

type UseCase = ucdef.UserAction[*Request, []actionlog.ActionLog]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-action-logs" }

func (uc *usecase) Execute(ctx context.Context, in *Request) ([]actionlog.ActionLog, error) {
	// Parse and validate time range (from < to)
	from, err := time.Parse(time.RFC3339, in.From)
	if err != nil {
		return nil, errx.New(
			"invalid 'from' time format, expected RFC3339",
			errx.WithType(errx.T_Validation),
			errx.WithCode("INVALID_TIME_FORMAT"),
		)
	}

	to, err := time.Parse(time.RFC3339, in.To)
	if err != nil {
		return nil, errx.New(
			"invalid 'to' time format, expected RFC3339",
			errx.WithType(errx.T_Validation),
			errx.WithCode("INVALID_TIME_FORMAT"),
		)
	}

	if !from.Before(to) {
		return nil, errx.New(
			"'from' must be before 'to'",
			errx.WithType(errx.T_Validation),
			errx.WithCode("INVALID_TIME_RANGE"),
		)
	}

	// List action logs matching filter criteria
	logs, err := uc.domainContainer.ActionLogRepo().List(ctx, actionlog.Filter{
		Module:      in.Module,
		OperationID: in.OperationID,
		UserID:      in.UserID,
		Tags:        in.Tags,
		GroupKey:    in.GroupKey,
		CreatedFrom: &from,
		CreatedTo:   &to,
		Cursor:      in.Cursor,
		Limit:       in.Limit,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return action logs
	return logs, nil
}
