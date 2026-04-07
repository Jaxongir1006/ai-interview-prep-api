package getuserpermissions

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
	Content []rbac.UserPermission `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-user-permissions" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find user by ID
	_, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.UserID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// List direct permissions for the user
	perms, err := uc.domainContainer.UserPermissionRepo().List(ctx, rbac.UserPermissionFilter{
		UserID: &in.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{Content: perms}, nil
}
