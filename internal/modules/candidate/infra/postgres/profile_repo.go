package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewProfileRepo(idb bun.IDB) profile.Repo {
	return repogen.NewPgRepoBuilder[profile.CandidateProfile, profile.Filter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(profile.CodeCandidateProfileNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_candidate_profiles_user_id": profile.CodeCandidateProfileUserConflict,
		}).
		WithFilterFunc(profileFilterFunc).
		Build()
}

func profileFilterFunc(q *bun.SelectQuery, f profile.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.TargetRole != nil {
		q = q.Where("target_role = ?", *f.TargetRole)
	}
	if f.ExperienceLevel != nil {
		q = q.Where("experience_level = ?", *f.ExperienceLevel)
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
