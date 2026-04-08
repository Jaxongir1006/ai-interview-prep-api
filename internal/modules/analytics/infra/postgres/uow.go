package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/achievement"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase/pguowbase"

	"github.com/uptrace/bun"
)

func NewUOWFactory(db *bun.DB) uow.Factory {
	return pguowbase.NewGenericFactory(
		db,
		schemaName,
		func(base *pguowbase.Base) uow.UnitOfWork {
			return &pgUOW{Base: base}
		},
	)
}

type pgUOW struct {
	*pguowbase.Base
}

func (u *pgUOW) ProgressSummary() progress.SummaryRepo {
	return NewProgressSummaryRepo(u.IDB())
}

func (u *pgUOW) TopicStat() progress.TopicStatRepo {
	return NewTopicStatRepo(u.IDB())
}

func (u *pgUOW) AchievementDefinition() achievement.DefinitionRepo {
	return NewAchievementDefinitionRepo(u.IDB())
}

func (u *pgUOW) CandidateAchievement() achievement.CandidateAchievementRepo {
	return NewCandidateAchievementRepo(u.IDB())
}
