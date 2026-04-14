package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/session"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewSessionRepo(idb bun.IDB) session.Repo {
	return repogen.NewPgRepoBuilder[session.Session, session.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(session.CodeSessionNotFound).
		WithConflictCodesMap(map[string]string{
			"interview_sessions_pkey": session.CodeSessionIDConflict,
		}).
		WithFilterFunc(sessionFilterFunc).
		Build()
}

func sessionFilterFunc(q *bun.SelectQuery, f session.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.TargetRole != nil {
		q = q.Where("target_role = ?", *f.TargetRole)
	}
	if f.ExperienceLevel != nil {
		q = q.Where("experience_level = ?", *f.ExperienceLevel)
	}
	if f.Difficulty != nil {
		q = q.Where("difficulty = ?", *f.Difficulty)
	}
	if f.StartedAtFrom != nil {
		q = q.Where("started_at >= ?", *f.StartedAtFrom)
	}
	if f.StartedAtTo != nil {
		q = q.Where("started_at <= ?", *f.StartedAtTo)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.Statuses != nil {
		q = q.Where("status IN (?)", bun.In(f.Statuses))
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
