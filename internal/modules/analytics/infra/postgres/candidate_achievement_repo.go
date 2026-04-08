package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/achievement"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewCandidateAchievementRepo(idb bun.IDB) achievement.CandidateAchievementRepo {
	return repogen.NewPgRepoBuilder[achievement.CandidateAchievement, achievement.CandidateAchievementFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(achievement.CodeCandidateAchievementNotFound).
		WithFilterFunc(candidateAchievementFilterFunc).
		Build()
}

func candidateAchievementFilterFunc(q *bun.SelectQuery, f achievement.CandidateAchievementFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.AchievementDefinitionID != nil {
		q = q.Where("achievement_definition_id = ?", *f.AchievementDefinitionID)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	return q
}
