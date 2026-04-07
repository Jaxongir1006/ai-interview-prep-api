package postgres

import (
	"context"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/audit/domain/statuschangelog"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

type statusChangeLogRepo struct {
	statuschangelog.Repo
	db *bun.DB
}

func NewStatusChangeLogRepo(idb bun.IDB, db *bun.DB) statuschangelog.Repo {
	inner := repogen.NewPgRepoBuilder[statuschangelog.StatusChangeLog, statuschangelog.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(statuschangelog.CodeStatusChangeLogNotFound).
		WithFilterFunc(statusChangeLogFilterFunc).
		Build()
	return &statusChangeLogRepo{Repo: inner, db: db}
}

func (r *statusChangeLogRepo) Create(
	ctx context.Context,
	entity *statuschangelog.StatusChangeLog,
) (*statuschangelog.StatusChangeLog, error) {
	return createWithPartitions(ctx, r.db, r.Repo.Create, entity)
}

func statusChangeLogFilterFunc(q *bun.SelectQuery, f statuschangelog.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.ActionLogID != nil {
		q = q.Where("action_log_id = ?", *f.ActionLogID)
	}
	if f.EntityType != nil {
		q = q.Where("entity_type = ?", *f.EntityType)
	}
	if f.EntityID != nil {
		q = q.Where("entity_id = ?", *f.EntityID)
	}
	if f.TraceID != nil {
		q = q.Where("trace_id = ?", *f.TraceID)
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
