package postgres

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/session"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/uow"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"
	"github.com/jaxongir1006/hire-ready-api/pkg/uowbase/pguowbase"

	"github.com/uptrace/bun"
)

func NewUOWFactory(db *bun.DB) uow.Factory {
	return pguowbase.NewGenericFactory(
		db,
		schemaName,
		func(base *pguowbase.Base) uow.UnitOfWork {
			return &pgUOW{Base: base}
		},
	)
}

type pgUOW struct {
	*pguowbase.Base
}

func (u *pgUOW) Role() rbac.RoleRepo {
	return NewRoleRepo(u.IDB())
}

func (u *pgUOW) RolePermission() rbac.RolePermissionRepo {
	return NewRolePermissionRepo(u.IDB())
}

func (u *pgUOW) UserRole() rbac.UserRoleRepo {
	return NewUserRoleRepo(u.IDB())
}

func (u *pgUOW) UserPermission() rbac.UserPermissionRepo {
	return NewUserPermissionRepo(u.IDB())
}

func (u *pgUOW) Session() session.Repo {
	return NewSessionRepo(u.IDB())
}

func (u *pgUOW) User() user.Repo {
	return NewUserRepo(u.IDB())
}
