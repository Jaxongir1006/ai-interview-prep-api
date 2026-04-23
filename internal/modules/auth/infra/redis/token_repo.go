package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/passwordresettoken"

	"github.com/code19m/errx"
	goredis "github.com/redis/go-redis/v9"
)

const (
	emailVerificationKeyPrefix = "auth:email_verification:"
	passwordResetKeyPrefix     = "auth:password_reset:"
)

type EmailVerificationTokenRepo struct {
	client goredis.Cmdable
}

func NewEmailVerificationTokenRepo(client goredis.Cmdable) emailverificationtoken.Repo {
	return &EmailVerificationTokenRepo{client: client}
}

func (r *EmailVerificationTokenRepo) Create(
	ctx context.Context,
	tokenHash string,
	token *emailverificationtoken.EmailVerificationToken,
	ttl time.Duration,
) error {
	err := createToken(
		ctx,
		r.client,
		emailVerificationKeyPrefix,
		tokenHash,
		token,
		ttl,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func (r *EmailVerificationTokenRepo) Consume(
	ctx context.Context,
	tokenHash string,
) (*emailverificationtoken.EmailVerificationToken, error) {
	token := &emailverificationtoken.EmailVerificationToken{}
	err := consumeToken(ctx, r.client, emailVerificationKeyPrefix, tokenHash, token)
	if errx.IsCodeIn(err, emailverificationtoken.CodeEmailVerificationTokenNotFound) {
		return nil, errx.Wrap(err)
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return token, nil
}

func (r *EmailVerificationTokenRepo) InvalidateUserEmail(ctx context.Context, userID, email string) error {
	err := invalidateUserEmail(ctx, r.client, emailVerificationKeyPrefix, userID, email)
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

type PasswordResetTokenRepo struct {
	client goredis.Cmdable
}

func NewPasswordResetTokenRepo(client goredis.Cmdable) passwordresettoken.Repo {
	return &PasswordResetTokenRepo{client: client}
}

func (r *PasswordResetTokenRepo) Create(
	ctx context.Context,
	tokenHash string,
	token *passwordresettoken.PasswordResetToken,
	ttl time.Duration,
) error {
	err := createToken(
		ctx,
		r.client,
		passwordResetKeyPrefix,
		tokenHash,
		token,
		ttl,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func (r *PasswordResetTokenRepo) Consume(
	ctx context.Context,
	tokenHash string,
) (*passwordresettoken.PasswordResetToken, error) {
	token := &passwordresettoken.PasswordResetToken{}
	err := consumeToken(ctx, r.client, passwordResetKeyPrefix, tokenHash, token)
	if errx.IsCodeIn(err, passwordresettoken.CodePasswordResetTokenNotFound) {
		return nil, errx.Wrap(err)
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return token, nil
}

func (r *PasswordResetTokenRepo) InvalidateUserEmail(ctx context.Context, userID, email string) error {
	err := invalidateUserEmail(ctx, r.client, passwordResetKeyPrefix, userID, email)
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func createToken(
	ctx context.Context,
	client goredis.Cmdable,
	keyPrefix,
	tokenHash string,
	token any,
	ttl time.Duration,
) error {
	userID, email, err := tokenIdentity(token)
	if err != nil {
		return errx.Wrap(err)
	}

	err = invalidateUserEmail(ctx, client, keyPrefix, userID, email)
	if err != nil {
		return errx.Wrap(err)
	}

	payload, err := json.Marshal(token)
	if err != nil {
		return errx.Wrap(err)
	}

	err = client.Set(ctx, tokenKey(keyPrefix, tokenHash), payload, ttl).Err()
	if err != nil {
		return errx.Wrap(err)
	}

	err = client.Set(ctx, userKey(keyPrefix, userID, email), tokenHash, ttl).Err()
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func consumeToken(
	ctx context.Context,
	client goredis.Cmdable,
	keyPrefix,
	tokenHash string,
	dest any,
) error {
	payload, err := client.GetDel(ctx, tokenKey(keyPrefix, tokenHash)).Bytes()
	if errors.Is(err, goredis.Nil) {
		return tokenNotFoundErr(keyPrefix)
	}
	if err != nil {
		return errx.Wrap(err)
	}

	err = json.Unmarshal(payload, dest)
	if err != nil {
		return errx.Wrap(err)
	}

	userID, email, err := tokenIdentity(dest)
	if err != nil {
		return errx.Wrap(err)
	}

	err = client.Del(ctx, userKey(keyPrefix, userID, email)).Err()
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func invalidateUserEmail(
	ctx context.Context,
	client goredis.Cmdable,
	keyPrefix,
	userID,
	email string,
) error {
	pointerKey := userKey(keyPrefix, userID, email)
	oldTokenHash, err := client.Get(ctx, pointerKey).Result()
	if err != nil && !errors.Is(err, goredis.Nil) {
		return errx.Wrap(err)
	}

	keys := []string{pointerKey}
	if oldTokenHash != "" {
		keys = append(keys, tokenKey(keyPrefix, oldTokenHash))
	}

	err = client.Del(ctx, keys...).Err()
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func tokenIdentity(token any) (string, string, error) {
	switch v := token.(type) {
	case *emailverificationtoken.EmailVerificationToken:
		return v.UserID, v.Email, nil
	case *passwordresettoken.PasswordResetToken:
		return v.UserID, v.Email, nil
	default:
		return "", "", errx.New("unsupported auth token type")
	}
}

func tokenNotFoundErr(keyPrefix string) error {
	if keyPrefix == emailVerificationKeyPrefix {
		return errx.New(
			"email verification token is not found",
			errx.WithCode(emailverificationtoken.CodeEmailVerificationTokenNotFound),
		)
	}

	return errx.New(
		"password reset token is not found",
		errx.WithCode(passwordresettoken.CodePasswordResetTokenNotFound),
	)
}

func tokenKey(prefix, tokenHash string) string {
	return prefix + "token:" + tokenHash
}

func userKey(prefix, userID, email string) string {
	return prefix + "user:" + userID + ":" + email
}
