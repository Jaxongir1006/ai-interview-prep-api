package getrecommendations

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/pblc/dashboard"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct{}

type UseCase = ucdef.UserAction[*Request, *dashboard.Recommendations]

func New(builder *dashboard.Builder) UseCase {
	return &usecase{builder: builder}
}

type usecase struct {
	builder *dashboard.Builder
}

func (uc *usecase) OperationID() string { return "get-dashboard-recommendations" }

func (uc *usecase) Execute(ctx context.Context, _ *Request) (*dashboard.Recommendations, error) {
	// Read authenticated user context
	userCtx := auth.MustUserContext(ctx)

	// Return dashboard recommendations
	out, err := uc.builder.Recommendations(ctx, userCtx.UserID)
	if err != nil {
		return nil, dashboardErr(err)
	}
	return &out, nil
}

func dashboardErr(err error) error {
	return errx.New(
		"dashboard data is temporarily unavailable",
		errx.WithType(errx.T_Internal),
		errx.WithCode(dashboard.CodeDashboardUnavailable),
		errx.WithDetails(errx.D{"cause": err.Error()}),
	)
}
