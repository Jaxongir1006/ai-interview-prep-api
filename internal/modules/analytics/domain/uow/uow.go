package uow

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/achievement"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"
)

type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
	uowbase.UnitOfWork

	ProgressSummary() progress.SummaryRepo
	TopicStat() progress.TopicStatRepo
	AchievementDefinition() achievement.DefinitionRepo
	CandidateAchievement() achievement.CandidateAchievementRepo
}
