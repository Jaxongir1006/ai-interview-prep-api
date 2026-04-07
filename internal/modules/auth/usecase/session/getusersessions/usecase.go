package getusersessions

import (
	"context"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/session"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/pagination"
	"github.com/rise-and-shine/pkg/sorter"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

type Request struct {
	pagination.Request

	UserID string `query:"user_id" validate:"required"`
	Sort   string `query:"sort"`
}

type UseCase = ucdef.UserAction[*Request, *pagination.Response[session.Session]]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-user-sessions" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[session.Session], error) {
	// Normalize pagination params
	in.Normalize()

	// Find user by ID
	_, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.UserID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// List sessions for the user
	sessions, count, err := uc.domainContainer.SessionRepo().ListWithCount(ctx, session.Filter{
		UserID:   &in.UserID,
		Limit:    lo.ToPtr(in.Limit()),
		Offset:   lo.ToPtr(in.Offset()),
		SortOpts: sorter.MakeFromStr(in.Sort, "last_used_at", "created_at"),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	resp := pagination.NewResponse(sessions, int64(count), in.Request)
	return &resp, nil
}
