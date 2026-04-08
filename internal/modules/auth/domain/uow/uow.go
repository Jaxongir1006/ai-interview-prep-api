package uow

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/oauthaccount"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/rbac"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"
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
	OAuthAccount() oauthaccount.Repo
	Session() session.Repo
	User() user.Repo
}
