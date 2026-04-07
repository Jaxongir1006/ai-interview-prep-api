package getusers

import (
	"context"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/pagination"
	"github.com/rise-and-shine/pkg/sorter"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

type Request struct {
	pagination.Request

	ID       *string `query:"id"`
	Username *string `query:"username"`
	IsActive *bool   `query:"is_active"`
	Sort     string  `query:"sort"`
}

type UserItem struct {
	ID                string     `json:"id"`
	Username          *string    `json:"username"`
	IsActive          bool       `json:"is_active"`
	LastActiveAt      *time.Time `json:"last_active_at"`
	Roles             []string   `json:"roles"`
	DirectPermissions []string   `json:"direct_permissions"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type UseCase = ucdef.UserAction[*Request, *pagination.Response[UserItem]]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-users" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[UserItem], error) {
	// Normalize pagination params
	in.Normalize()

	// List users matching filter criteria
	users, count, err := uc.domainContainer.UserRepo().ListWithCount(ctx, user.Filter{
		ID:       in.ID,
		Username: in.Username,
		IsActive: in.IsActive,
		Limit:    lo.ToPtr(in.Limit()),
		Offset:   lo.ToPtr(in.Offset()),
		SortOpts: sorter.MakeFromStr(in.Sort, "username", "created_at", "updated_at"),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(users) == 0 {
		resp := pagination.NewResponse([]UserItem{}, int64(count), in.Request)
		return &resp, nil
	}

	// Collect all user IDs
	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	// Batch-fetch role assignments for all users
	userRoles, err := uc.domainContainer.UserRoleRepo().List(ctx, rbac.UserRoleFilter{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Collect unique role IDs and fetch role names
	roleIDSet := make(map[int64]struct{})
	for _, ur := range userRoles {
		roleIDSet[ur.RoleID] = struct{}{}
	}

	roleNameMap := make(map[int64]string)
	if len(roleIDSet) > 0 {
		roleIDs := make([]int64, 0, len(roleIDSet))
		for id := range roleIDSet {
			roleIDs = append(roleIDs, id)
		}

		roles, rolesErr := uc.domainContainer.RoleRepo().List(ctx, rbac.RoleFilter{
			IDs: roleIDs,
		})
		if rolesErr != nil {
			return nil, errx.Wrap(rolesErr)
		}

		for _, r := range roles {
			roleNameMap[r.ID] = r.Name
		}
	}

	// Build user → role names map
	userRolesMap := make(map[string][]string)
	for _, ur := range userRoles {
		if name, ok := roleNameMap[ur.RoleID]; ok {
			userRolesMap[ur.UserID] = append(userRolesMap[ur.UserID], name)
		}
	}

	// Batch-fetch direct permissions for all users
	userPerms, err := uc.domainContainer.UserPermissionRepo().List(ctx, rbac.UserPermissionFilter{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Build user → permission strings map
	userPermsMap := make(map[string][]string)
	for _, up := range userPerms {
		userPermsMap[up.UserID] = append(userPermsMap[up.UserID], up.Permission)
	}

	// Assemble response
	items := make([]UserItem, len(users))
	for i, u := range users {
		roles := userRolesMap[u.ID]
		if roles == nil {
			roles = []string{}
		}
		perms := userPermsMap[u.ID]
		if perms == nil {
			perms = []string{}
		}

		items[i] = UserItem{
			ID:                u.ID,
			Username:          u.Username,
			IsActive:          u.IsActive,
			LastActiveAt:      u.LastActiveAt,
			Roles:             roles,
			DirectPermissions: perms,
			CreatedAt:         u.CreatedAt,
			UpdatedAt:         u.UpdatedAt,
		}
	}

	resp := pagination.NewResponse(items, int64(count), in.Request)
	return &resp, nil
}
