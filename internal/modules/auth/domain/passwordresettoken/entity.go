package passwordresettoken

import (
	"context"
	"time"
)

const (
	CodePasswordResetTokenNotFound = "PASSWORD_RESET_TOKEN_NOT_FOUND"
	CodePasswordResetTokenInvalid  = "PASSWORD_RESET_TOKEN_INVALID"
)

type PasswordResetToken struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Repo interface {
	Create(ctx context.Context, tokenHash string, token *PasswordResetToken, ttl time.Duration) error
	Consume(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	InvalidateUserEmail(ctx context.Context, userID, email string) error
}
