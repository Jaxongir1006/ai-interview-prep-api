package postgres

import (
	"context"

	"go-enterprise-blueprint/internal/modules/audit/domain/actionlog"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type actionLogRepo struct {
	actionlog.Repo
	db *bun.DB
}

func NewActionLogRepo(idb bun.IDB, db *bun.DB) actionlog.Repo {
	inner := repogen.NewPgRepoBuilder[actionlog.ActionLog, actionlog.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(actionlog.CodeActionLogNotFound).
		WithFilterFunc(actionLogFilterFunc).
		Build()
	return &actionLogRepo{Repo: inner, db: db}
}

func (r *actionLogRepo) Create(ctx context.Context, entity *actionlog.ActionLog) (*actionlog.ActionLog, error) {
	return createWithPartitions(ctx, r.db, r.Repo.Create, entity)
}

func actionLogFilterFunc(q *bun.SelectQuery, f actionlog.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.Module != nil {
		q = q.Where("module = ?", *f.Module)
	}
	if f.OperationID != nil {
		q = q.Where("operation_id = ?", *f.OperationID)
	}
	if f.TraceID != nil {
		q = q.Where("trace_id = ?", *f.TraceID)
	}
	if f.Tags != nil {
		q = q.Where("tags @> ?::varchar[]", pgdialect.Array(f.Tags))
	}
	if f.GroupKey != nil {
		q = q.Where("group_key = ?", *f.GroupKey)
	}
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at < ?", *f.CreatedTo)
	}
	if f.Cursor != nil {
		q = q.Where("id < ?", *f.Cursor)
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	q = q.Order("id DESC")
	return q
}
