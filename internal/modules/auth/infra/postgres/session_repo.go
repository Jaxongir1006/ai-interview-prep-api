package postgres

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/spf13/cast"
	"github.com/uptrace/bun"
)

type sessionRepo struct {
	repogen.Repo[session.Session, session.Filter]
	idb bun.IDB
}

func NewSessionRepo(idb bun.IDB) session.Repo {
	baseRepo := repogen.NewPgRepoBuilder[session.Session, session.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(session.CodeSessionNotFound).
		WithFilterFunc(sessionFilterFunc).
		Build()

	return &sessionRepo{
		Repo: baseRepo,
		idb:  idb,
	}
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) (int64, error) {
	now := time.Now()

	res, err := r.idb.NewDelete().
		Model((*session.Session)(nil)).
		ModelTableExpr(schemaName+".sessions AS s").
		Where("refresh_token_expires_at < ?", now).
		Exec(ctx)
	if err != nil {
		return 0, errx.Wrap(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, errx.Wrap(err)
	}

	return count, nil
}

func sessionFilterFunc(q *bun.SelectQuery, f session.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.AccessToken != nil {
		q = q.Where("access_token = ?", *f.AccessToken)
	}
	if f.RefreshToken != nil {
		q = q.Where("refresh_token = ?", *f.RefreshToken)
	}
	if f.IsActive != nil {
		if *f.IsActive {
			q = q.Where("refresh_token_expires_at >= ?", time.Now())
		} else {
			q = q.Where("refresh_token_expires_at < ?", time.Now())
		}
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	if f.Offset != nil {
		q = q.Offset(*f.Offset)
	}
	if f.OrderByLastUsedAt != nil {
		q = q.Order("last_used_at " + cast.ToString(f.OrderByLastUsedAt))
	}
	for _, o := range f.SortOpts {
		q = q.Order(o.ToSQL())
	}
	return q
}
