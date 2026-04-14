package getrecentactivity

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/pblc/dashboard"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/rise-and-shine/pkg/val"
)

type Request struct {
	Limit  int     `query:"limit"`
	Cursor *string `query:"cursor"`
}

type UseCase = ucdef.UserAction[*Request, *dashboard.RecentActivity]

func New(builder *dashboard.Builder) UseCase {
	return &usecase{builder: builder}
}

type usecase struct {
	builder *dashboard.Builder
}

func (uc *usecase) OperationID() string { return "get-recent-activity" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*dashboard.RecentActivity, error) {
	// Validate input
	limit := in.Limit
	if limit == 0 {
		limit = 10
	}
	if limit < 1 || limit > 50 {
		return nil, validationErr("limit")
	}

	// Read authenticated user context
	userCtx := auth.MustUserContext(ctx)

	// Return recent activity items
	out, err := uc.builder.RecentActivity(ctx, userCtx.UserID, limit, in.Cursor)
	if err != nil {
		return nil, dashboardErr(err)
	}
	return &out, nil
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
