package gettopics

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/pblc/dashboard"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/rise-and-shine/pkg/val"
)

type Request struct {
	Range string `query:"range"`
}

type Response struct {
	Range string                       `json:"range"`
	Items []dashboard.TopicPerformance `json:"items"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(builder *dashboard.Builder) UseCase {
	return &usecase{builder: builder}
}

type usecase struct {
	builder *dashboard.Builder
}

func (uc *usecase) OperationID() string { return "get-dashboard-topics" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Validate input
	rangeValue := dashboard.NormalizeRange(in.Range)
	if !dashboard.IsValidRange(rangeValue) {
		return nil, validationErr("range")
	}

	// Read authenticated user context
	userCtx := auth.MustUserContext(ctx)

	// Return topic performance items
	topics, err := uc.builder.Topics(ctx, userCtx.UserID, rangeValue)
	if err != nil {
		return nil, dashboardErr(err)
	}
	return &Response{Range: rangeValue, Items: topics.Items}, nil
}

func validationErr(field string) error {
	return errx.New(
		"dashboard request validation failed",
		errx.WithType(errx.T_Validation),
		errx.WithCode(val.CodeValidationFailed),
		errx.WithDetails(errx.D{"field": field}),
	)
}

func dashboardErr(err error) error {
	return errx.New(
		"dashboard data is temporarily unavailable",
		errx.WithType(errx.T_Internal),
		errx.WithCode(dashboard.CodeDashboardUnavailable),
		errx.WithDetails(errx.D{"cause": err.Error()}),
	)
}
