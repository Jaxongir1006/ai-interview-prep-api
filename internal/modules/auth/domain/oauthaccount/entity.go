package oauthaccount

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeOAuthAccountNotFound      = "OAUTH_ACCOUNT_NOT_FOUND"
	CodeOAuthProviderUserConflict = "OAUTH_PROVIDER_USER_CONFLICT"
	CodeOAuthUserProviderConflict = "OAUTH_USER_PROVIDER_CONFLICT"

	ProviderGoogle = "google"
	ProviderGitHub = "github"
)

type OAuthAccount struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	ProviderEmail  *string    `json:"provider_email"`
	LastLoginAt    *time.Time `json:"last_login_at"`
}
