package postgres

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewUserRoleRepo(idb bun.IDB) rbac.UserRoleRepo {
	return repogen.NewPgRepoBuilder[rbac.UserRole, rbac.UserRoleFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(rbac.CodeUserRoleNotFound).
		WithFilterFunc(userRoleFilterFunc).
		Build()
}

func userRoleFilterFunc(q *bun.SelectQuery, f rbac.UserRoleFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.RoleID != nil {
		q = q.Where("role_id = ?", *f.RoleID)
	}
	if f.UserIDs != nil {
		q = q.Where("user_id IN (?)", bun.In(f.UserIDs))
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	if f.Offset != nil {
		q = q.Offset(*f.Offset)
	}
	for _, o := range f.SortOpts {
		q = q.Order(o.ToSQL())
	}
	return q
}
