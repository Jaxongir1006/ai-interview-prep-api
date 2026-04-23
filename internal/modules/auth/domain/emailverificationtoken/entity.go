package emailverificationtoken

import (
	"context"
	"time"
)

const (
	CodeEmailVerificationTokenNotFound = "EMAIL_VERIFICATION_TOKEN_NOT_FOUND" //nolint:gosec // Public error code.
	CodeEmailVerificationTokenConflict = "EMAIL_VERIFICATION_TOKEN_CONFLICT"  //nolint:gosec // Public error code.
	CodeEmailVerificationTokenInvalid  = "EMAIL_VERIFICATION_TOKEN_INVALID"   //nolint:gosec // Public error code.
	CodeEmailVerificationTokenExpired  = "EMAIL_VERIFICATION_TOKEN_EXPIRED"   //nolint:gosec // Public error code.
	CodeEmailVerificationEmailMismatch = "EMAIL_VERIFICATION_EMAIL_MISMATCH"
)

type EmailVerificationToken struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Repo interface {
	Create(ctx context.Context, tokenHash string, token *EmailVerificationToken, ttl time.Duration) error
	Consume(ctx context.Context, tokenHash string) (*EmailVerificationToken, error)
	InvalidateUserEmail(ctx context.Context, userID, email string) error
}
