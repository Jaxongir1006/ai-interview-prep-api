package rbac

import (
	"context"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type RoleFilter struct {
	ID   *int64
	Name *string
	IDs  []int64

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type RolePermissionFilter struct {
	ID     *int64
	RoleID *int64

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type UserRoleFilter struct {
	ID     *int64
	UserID *string
	RoleID *int64

	UserIDs []string

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type UserPermissionFilter struct {
	ID     *int64
	UserID *string

	UserIDs []string

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type RoleRepo interface {
	repogen.Repo[Role, RoleFilter]
}

type RolePermissionRepo interface {
	repogen.Repo[RolePermission, RolePermissionFilter]
}

type UserRoleRepo interface {
	repogen.Repo[UserRole, UserRoleFilter]
}

type UserPermissionRepo interface {
	repogen.Repo[UserPermission, UserPermissionFilter]

	CollectUserPermissions(ctx context.Context, userID string) ([]string, error)
}
