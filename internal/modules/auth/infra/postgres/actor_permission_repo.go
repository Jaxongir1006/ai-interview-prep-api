package postgres

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

type userPermissionRepo struct {
	repogen.Repo[rbac.UserPermission, rbac.UserPermissionFilter]
	idb bun.IDB
}

func NewUserPermissionRepo(idb bun.IDB) rbac.UserPermissionRepo {
	return &userPermissionRepo{
		Repo: repogen.NewPgRepoBuilder[rbac.UserPermission, rbac.UserPermissionFilter](idb).
			WithSchemaName(schemaName).
			WithNotFoundCode(rbac.CodeUserPermissionNotFound).
			WithFilterFunc(userPermissionFilterFunc).
			Build(),
		idb: idb,
	}
}

func (r *userPermissionRepo) CollectUserPermissions(ctx context.Context, userID string) ([]string, error) {
	var permissions []string

	err := r.idb.NewRaw(`
		SELECT rp.permission FROM auth.role_permissions rp
		JOIN auth.user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?
		UNION
		SELECT up.permission FROM auth.user_permissions up
		WHERE up.user_id = ?
	`, userID, userID).Scan(ctx, &permissions)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return permissions, nil
}

func userPermissionFilterFunc(q *bun.SelectQuery, f rbac.UserPermissionFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
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
