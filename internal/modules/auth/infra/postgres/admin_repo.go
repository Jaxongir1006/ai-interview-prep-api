package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewUserRepo(idb bun.IDB) user.Repo {
	return repogen.NewPgRepoBuilder[user.User, user.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(user.CodeUserNotFound).
		WithConflictCodesMap(map[string]string{
			"users_username_key":    user.CodeUsernameConflict,
			"uk_users_email":        user.CodeEmailConflict,
			"uk_users_phone_number": user.CodePhoneNumberConflict,
		}).
		WithFilterFunc(userFilterFunc).
		Build()
}

func userFilterFunc(q *bun.SelectQuery, f user.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Username != nil {
		q = q.Where("username = ?", *f.Username)
	}
	if f.Email != nil {
		q = q.Where("email = ?", *f.Email)
	}
	if f.PhoneNumber != nil {
		q = q.Where("phone_number = ?", *f.PhoneNumber)
	}
	if f.IsVerified != nil {
		q = q.Where("is_verified = ?", *f.IsVerified)
	}
	if f.IsActive != nil {
		q = q.Where("is_active = ?", *f.IsActive)
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
