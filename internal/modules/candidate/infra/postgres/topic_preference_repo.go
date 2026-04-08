package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewTopicPreferenceRepo(idb bun.IDB) topicpreference.Repo {
	return repogen.NewPgRepoBuilder[topicpreference.TopicPreference, topicpreference.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(topicpreference.CodeTopicPreferenceNotFound).
		WithFilterFunc(topicPreferenceFilterFunc).
		Build()
}

func topicPreferenceFilterFunc(q *bun.SelectQuery, f topicpreference.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.CandidateProfileID != nil {
		q = q.Where("candidate_profile_id = ?", *f.CandidateProfileID)
	}
	if f.TopicKey != nil {
		q = q.Where("topic_key = ?", *f.TopicKey)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	if f.Offset != nil {
		q = q.Offset(*f.Offset)
	}
	for _, o := range f.SortOpts {
		q = q.Order(o.ToSQL())
	}
	return q
}
