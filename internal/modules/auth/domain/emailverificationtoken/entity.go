package emailverificationtoken

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeEmailVerificationTokenNotFound = "EMAIL_VERIFICATION_TOKEN_NOT_FOUND" //nolint:gosec // Public error code.
	CodeEmailVerificationTokenConflict = "EMAIL_VERIFICATION_TOKEN_CONFLICT"  //nolint:gosec // Public error code.
	CodeEmailVerificationTokenInvalid  = "EMAIL_VERIFICATION_TOKEN_INVALID"   //nolint:gosec // Public error code.
	CodeEmailVerificationTokenExpired  = "EMAIL_VERIFICATION_TOKEN_EXPIRED"   //nolint:gosec // Public error code.
	CodeEmailVerificationEmailMismatch = "EMAIL_VERIFICATION_EMAIL_MISMATCH"
)

type EmailVerificationToken struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`
	Email  string `json:"email"`

	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
}
