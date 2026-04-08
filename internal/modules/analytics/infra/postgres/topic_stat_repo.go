package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewTopicStatRepo(idb bun.IDB) progress.TopicStatRepo {
	return repogen.NewPgRepoBuilder[progress.TopicStat, progress.TopicStatFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(progress.CodeTopicStatNotFound).
		WithFilterFunc(topicStatFilterFunc).
		Build()
}

func topicStatFilterFunc(q *bun.SelectQuery, f progress.TopicStatFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.TopicKey != nil {
		q = q.Where("topic_key = ?", *f.TopicKey)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	return q
}
