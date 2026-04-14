package getperformancetrend

import (
	"context"
	"strings"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/pblc/dashboard"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/rise-and-shine/pkg/val"
)

type Request struct {
	Range   string  `query:"range"`
	TopicID *string `query:"topic_id"`
}

type Response struct {
	Range  string                       `json:"range"`
	Topic  *dashboard.Option            `json:"topic"`
	Points []dashboard.PerformancePoint `json:"points"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(builder *dashboard.Builder) UseCase {
	return &usecase{builder: builder}
}

type usecase struct {
	builder *dashboard.Builder
}

func (uc *usecase) OperationID() string { return "get-performance-trend" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Validate input
	rangeValue := dashboard.NormalizeRange(in.Range)
	if !dashboard.IsValidRange(rangeValue) {
		return nil, validationErr("range")
	}
	if in.TopicID != nil && strings.TrimSpace(*in.TopicID) == "" {
		return nil, validationErr("topic_id")
	}

	// Read authenticated user context
	userCtx := auth.MustUserContext(ctx)

	// Return performance trend points
	performance, err := uc.builder.Performance(ctx, userCtx.UserID, rangeValue, in.TopicID)
	if err != nil {
		return nil, dashboardErr(err)
	}
	topic := uc.builder.TopicOption(in.TopicID)
	return &Response{Range: rangeValue, Topic: topic, Points: performance.Points}, nil
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
