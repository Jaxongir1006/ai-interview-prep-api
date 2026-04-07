package getroles

import (
	"context"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/pagination"
	"github.com/rise-and-shine/pkg/sorter"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

type Request struct {
	pagination.Request

	ID   *int64  `query:"id"`
	Name *string `query:"name"`
	Sort string  `query:"sort"`
}

type UseCase = ucdef.UserAction[*Request, *pagination.Response[rbac.Role]]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-roles" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[rbac.Role], error) {
	// Normalize pagination params
	in.Normalize()

	// List roles matching filter criteria
	roles, count, err := uc.domainContainer.RoleRepo().ListWithCount(ctx, rbac.RoleFilter{
		ID:       in.ID,
		Name:     in.Name,
		Limit:    lo.ToPtr(in.Limit()),
		Offset:   lo.ToPtr(in.Offset()),
		SortOpts: sorter.MakeFromStr(in.Sort, "name", "created_at", "updated_at"),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	resp := pagination.NewResponse(roles, int64(count), in.Request)
	return &resp, nil
}
