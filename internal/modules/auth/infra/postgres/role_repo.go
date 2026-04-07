package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewRoleRepo(idb bun.IDB) rbac.RoleRepo {
	return repogen.NewPgRepoBuilder[rbac.Role, rbac.RoleFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(rbac.CodeRoleNotFound).
		WithConflictCodesMap(map[string]string{
			"roles_name_key": rbac.CodeRoleNameConflict,
		}).
		WithFilterFunc(roleFilterFunc).
		Build()
}

func roleFilterFunc(q *bun.SelectQuery, f rbac.RoleFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Name != nil {
		q = q.Where("name = ?", *f.Name)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
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
