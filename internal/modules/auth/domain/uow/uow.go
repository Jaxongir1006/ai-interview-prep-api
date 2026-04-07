package uow

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/pkg/uowbase"
)

// Factory defines an interface for creating new instances of the UnitOfWork.
type Factory = uowbase.Factory[UnitOfWork]

// UnitOfWork represents a single unit of work, typically mapping to a database transaction.
// It provides access to various repositories and methods to finalize or discard changes.
type UnitOfWork interface {
	uowbase.UnitOfWork

	// Repository accessors
	Role() rbac.RoleRepo
	RolePermission() rbac.RolePermissionRepo
	UserRole() rbac.UserRoleRepo
	UserPermission() rbac.UserPermissionRepo
	Session() session.Repo
	User() user.Repo
}
