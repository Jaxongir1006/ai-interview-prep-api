package getrolepermissions

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	RoleID int64 `query:"role_id" validate:"required"`
}

type Response struct {
	Content []rbac.RolePermission `json:"content"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-role-permissions" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find role by ID
	_, err := uc.domainContainer.RoleRepo().Get(ctx, rbac.RoleFilter{
		ID: &in.RoleID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, rbac.CodeRoleNotFound)
	}

	// List permissions for the role
	perms, err := uc.domainContainer.RolePermissionRepo().List(ctx, rbac.RolePermissionFilter{
		RoleID: &in.RoleID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{Content: perms}, nil
}
