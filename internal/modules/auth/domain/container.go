package domain

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/uow"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
)

// Container holds domain interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct {
	userRepo           user.Repo
	sessionRepo        session.Repo
	roleRepo           rbac.RoleRepo
	rolePermissionRepo rbac.RolePermissionRepo
	userRoleRepo       rbac.UserRoleRepo
	userPermissionRepo rbac.UserPermissionRepo
	uowFactory         uow.Factory
}

func NewContainer(
	userRepo user.Repo,
	sessionRepo session.Repo,
	roleRepo rbac.RoleRepo,
	rolePermissionRepo rbac.RolePermissionRepo,
	userRoleRepo rbac.UserRoleRepo,
	userPermissionRepo rbac.UserPermissionRepo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		userRepo,
		sessionRepo,
		roleRepo,
		rolePermissionRepo,
		userRoleRepo,
		userPermissionRepo,
		uowFactory,
	}
}

func (c *Container) UserRepo() user.Repo {
	return c.userRepo
}

func (c *Container) SessionRepo() session.Repo {
	return c.sessionRepo
}

func (c *Container) RoleRepo() rbac.RoleRepo {
	return c.roleRepo
}

func (c *Container) RolePermissionRepo() rbac.RolePermissionRepo {
	return c.rolePermissionRepo
}

func (c *Container) UserRoleRepo() rbac.UserRoleRepo {
	return c.userRoleRepo
}

func (c *Container) UserPermissionRepo() rbac.UserPermissionRepo {
	return c.userPermissionRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
