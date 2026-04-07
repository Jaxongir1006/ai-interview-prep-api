package getuserroles

import (
	"context"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	UserID string `query:"user_id" validate:"required"`
}

type Response struct {
	Content []rbac.Role `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-user-roles" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find user by ID
	_, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.UserID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// List role assignments for the user
	userRoles, err := uc.domainContainer.UserRoleRepo().List(ctx, rbac.UserRoleFilter{
		UserID: &in.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(userRoles) == 0 {
		return &Response{Content: []rbac.Role{}}, nil
	}

	// Fetch role details for each assignment
	roleIDs := make([]int64, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	roles, err := uc.domainContainer.RoleRepo().List(ctx, rbac.RoleFilter{
		IDs: roleIDs,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{Content: roles}, nil
}
