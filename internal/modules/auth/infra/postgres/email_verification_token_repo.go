package postgres

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

type emailVerificationTokenRepo struct {
	repogen.Repo[emailverificationtoken.EmailVerificationToken, emailverificationtoken.Filter]
	idb bun.IDB
}

func NewEmailVerificationTokenRepo(idb bun.IDB) emailverificationtoken.Repo {
	baseRepo := repogen.NewPgRepoBuilder[
		emailverificationtoken.EmailVerificationToken,
		emailverificationtoken.Filter,
	](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(emailverificationtoken.CodeEmailVerificationTokenNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_email_verification_tokens_token_hash": emailverificationtoken.CodeEmailVerificationTokenConflict,
		}).
		WithFilterFunc(emailVerificationTokenFilterFunc).
		Build()

	return &emailVerificationTokenRepo{
		Repo: baseRepo,
		idb:  idb,
	}
}

func (r *emailVerificationTokenRepo) ExpireUnused(
	ctx context.Context,
	userID,
	email string,
	now time.Time,
) error {
	_, err := r.idb.NewUpdate().
		Model((*emailverificationtoken.EmailVerificationToken)(nil)).
		ModelTableExpr(schemaName+".email_verification_tokens AS evt").
		Set("expires_at = ?", now).
		Where("user_id = ?", userID).
		Where("email = ?", email).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		Exec(ctx)
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func emailVerificationTokenFilterFunc(
	q *bun.SelectQuery,
	f emailverificationtoken.Filter,
) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.Email != nil {
		q = q.Where("email = ?", *f.Email)
	}
	if f.TokenHash != nil {
		q = q.Where("token_hash = ?", *f.TokenHash)
	}
	if f.Unused != nil {
		if *f.Unused {
			q = q.Where("used_at IS NULL")
		} else {
			q = q.Where("used_at IS NOT NULL")
		}
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	return q
}
