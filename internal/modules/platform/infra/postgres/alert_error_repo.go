package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-enterprise-blueprint/internal/modules/platform/domain/alerterror"

	"github.com/code19m/errx"
	"github.com/uptrace/bun"
)

type alertErrorRepo struct {
	db     *bun.DB
	schema string
}

func NewAlertErrorRepo(db *bun.DB, schema string) alerterror.Repo {
	return &alertErrorRepo{db: db, schema: schema}
}

func (r *alertErrorRepo) table() string {
	return fmt.Sprintf("%s.errors", r.schema)
}

func (r *alertErrorRepo) Get(ctx context.Context, id string) (*alerterror.Error, error) {
	var e alerterror.Error

	err := r.db.NewSelect().
		TableExpr(r.table()).
		Where("id = ?", id).
		Scan(ctx, &e)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.New("error not found", errx.WithCode(alerterror.CodeErrorNotFound))
		}
		return nil, errx.Wrap(err)
	}

	return &e, nil
}

func (r *alertErrorRepo) ListWithCount(ctx context.Context, f alerterror.Filter) ([]alerterror.Error, int64, error) {
	q := r.db.NewSelect().
		TableExpr(r.table())

	q = r.applyFilter(q, f)

	count, err := q.Count(ctx)
	if err != nil {
		return nil, 0, errx.Wrap(err)
	}

	var items []alerterror.Error

	q = r.db.NewSelect().
		TableExpr(r.table())

	q = r.applyFilter(q, f)

	if f.Offset != nil {
		q = q.Offset(*f.Offset)
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}

	for _, o := range f.SortOpts {
		q = q.Order(o.ToSQL())
	}
	if len(f.SortOpts) == 0 {
		q = q.Order("created_at DESC")
	}

	err = q.Scan(ctx, &items)
	if err != nil {
		return nil, 0, errx.Wrap(err)
	}

	return items, int64(count), nil
}

func (r *alertErrorRepo) GetStats(ctx context.Context, f alerterror.StatsFilter) (*alerterror.Stats, error) {
	stats := &alerterror.Stats{}

	baseQ := r.db.NewSelect().TableExpr(r.table())
	baseQ = r.applyStatsFilter(baseQ, f)

	// Total count
	count, err := baseQ.Count(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	stats.TotalCount = int64(count)

	// By service
	byServiceQ := r.db.NewSelect().
		TableExpr(r.table()).
		ColumnExpr("service").
		ColumnExpr("COUNT(*) AS count")
	byServiceQ = r.applyStatsFilter(byServiceQ, f)
	byServiceQ = byServiceQ.GroupExpr("service").OrderExpr("count DESC")

	err = byServiceQ.Scan(ctx, &stats.ByService)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// By operation
	byOperationQ := r.db.NewSelect().
		TableExpr(r.table()).
		ColumnExpr("operation").
		ColumnExpr("COUNT(*) AS count")
	byOperationQ = r.applyStatsFilter(byOperationQ, f)
	byOperationQ = byOperationQ.GroupExpr("operation").OrderExpr("count DESC")

	err = byOperationQ.Scan(ctx, &stats.ByOperation)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// By code
	byCodeQ := r.db.NewSelect().
		TableExpr(r.table()).
		ColumnExpr("code").
		ColumnExpr("COUNT(*) AS count")
	byCodeQ = r.applyStatsFilter(byCodeQ, f)
	byCodeQ = byCodeQ.GroupExpr("code").OrderExpr("count DESC")

	err = byCodeQ.Scan(ctx, &stats.ByCode)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return stats, nil
}

func (r *alertErrorRepo) DeleteOlderThan(ctx context.Context, before time.Time) (int64, error) {
	res, err := r.db.NewDelete().
		TableExpr(r.table()).
		Where("created_at < ?", before).
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

func (r *alertErrorRepo) applyFilter(q *bun.SelectQuery, f alerterror.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Code != nil {
		q = q.Where("code = ?", *f.Code)
	}
	if f.Service != nil {
		q = q.Where("service = ?", *f.Service)
	}
	if f.Operation != nil {
		q = q.Where("operation = ?", *f.Operation)
	}
	if f.Alerted != nil {
		q = q.Where("alerted = ?", *f.Alerted)
	}
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at < ?", *f.CreatedTo)
	}
	if f.Search != "" {
		search := "%" + f.Search + "%"
		q = q.WhereGroup(" AND ", func(wq *bun.SelectQuery) *bun.SelectQuery {
			return wq.
				Where("code ILIKE ?", search).
				WhereOr("message ILIKE ?", search).
				WhereOr("operation ILIKE ?", search)
		})
	}
	return q
}

func (r *alertErrorRepo) applyStatsFilter(q *bun.SelectQuery, f alerterror.StatsFilter) *bun.SelectQuery {
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at < ?", *f.CreatedTo)
	}
	return q
}
