package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewProgressSummaryRepo(idb bun.IDB) progress.SummaryRepo {
	return repogen.NewPgRepoBuilder[progress.Summary, progress.SummaryFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(progress.CodeProgressSummaryNotFound).
		WithFilterFunc(progressSummaryFilterFunc).
		Build()
}

func progressSummaryFilterFunc(q *bun.SelectQuery, f progress.SummaryFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	return q
}
