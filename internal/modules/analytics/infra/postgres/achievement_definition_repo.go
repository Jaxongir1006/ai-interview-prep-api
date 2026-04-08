package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/achievement"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewAchievementDefinitionRepo(idb bun.IDB) achievement.DefinitionRepo {
	return repogen.NewPgRepoBuilder[achievement.Definition, achievement.DefinitionFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(achievement.CodeAchievementDefinitionNotFound).
		WithConflictCodesMap(map[string]string{
			"uk_achievement_definitions_code": achievement.CodeAchievementCodeConflict,
		}).
		WithFilterFunc(achievementDefinitionFilterFunc).
		Build()
}

func achievementDefinitionFilterFunc(q *bun.SelectQuery, f achievement.DefinitionFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Code != nil {
		q = q.Where("code = ?", *f.Code)
	}
	if f.Category != nil {
		q = q.Where("category = ?", *f.Category)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	return q
}
