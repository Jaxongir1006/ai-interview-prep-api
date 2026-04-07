package rbac

import (
	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeRoleNotFound           = "ROLE_NOT_FOUND"
	CodeRoleNameConflict       = "ROLE_NAME_CONFLICT"
	CodeRolePermissionNotFound = "ROLE_PERMISSION_NOT_FOUND"
	CodeUserRoleNotFound       = "USER_ROLE_NOT_FOUND"
	CodeUserPermissionNotFound = "USER_PERMISSION_NOT_FOUND"
	CodeRoleHasAssignedUsers   = "ROLE_HAS_ASSIGNED_USERS"
)

type Role struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	// Name is a unique name of the role
	Name string `json:"name"`
}

type RolePermission struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	RoleID     int64  `json:"role_id"`
	Permission string `json:"permission"`
}

type UserRole struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`
	RoleID int64  `json:"role_id"`
}

type UserPermission struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
}
