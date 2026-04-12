package oauthaccount

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeOAuthAccountNotFound      = "OAUTH_ACCOUNT_NOT_FOUND"
	CodeOAuthProviderUserConflict = "OAUTH_PROVIDER_USER_CONFLICT"
	CodeOAuthUserProviderConflict = "OAUTH_USER_PROVIDER_CONFLICT"

	ProviderGoogle = "google"
	ProviderGitHub = "github"
)

type OAuthAccount struct {
	bun.BaseModel `bun:"table:oauth_accounts,alias:oa"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	ProviderEmail  *string    `json:"provider_email"`
	LastLoginAt    *time.Time `json:"last_login_at"`

	CreatedAt time.Time `bun:",nullzero" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero" json:"updated_at"`
}

func (m *OAuthAccount) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
