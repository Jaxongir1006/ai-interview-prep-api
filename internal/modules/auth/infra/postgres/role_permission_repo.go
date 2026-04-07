package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewRolePermissionRepo(idb bun.IDB) rbac.RolePermissionRepo {
	return repogen.NewPgRepoBuilder[rbac.RolePermission, rbac.RolePermissionFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(rbac.CodeRolePermissionNotFound).
		WithFilterFunc(rolePermissionFilterFunc).
		Build()
}

func rolePermissionFilterFunc(q *bun.SelectQuery, f rbac.RolePermissionFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.RoleID != nil {
		q = q.Where("role_id = ?", *f.RoleID)
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
