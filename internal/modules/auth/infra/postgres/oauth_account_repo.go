package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/oauthaccount"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewOAuthAccountRepo(idb bun.IDB) oauthaccount.Repo {
	return repogen.NewPgRepoBuilder[oauthaccount.OAuthAccount, oauthaccount.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(oauthaccount.CodeOAuthAccountNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_oauth_accounts_provider_user": oauthaccount.CodeOAuthProviderUserConflict,
			"uk_oauth_accounts_user_provider": oauthaccount.CodeOAuthUserProviderConflict,
		}).
		WithFilterFunc(oauthAccountFilterFunc).
		Build()
}

func oauthAccountFilterFunc(q *bun.SelectQuery, f oauthaccount.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.Provider != nil {
		q = q.Where("provider = ?", *f.Provider)
	}
	if f.ProviderUserID != nil {
		q = q.Where("provider_user_id = ?", *f.ProviderUserID)
	}
	if f.ProviderEmail != nil {
		q = q.Where("provider_email = ?", *f.ProviderEmail)
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
