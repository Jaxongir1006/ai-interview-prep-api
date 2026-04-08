package domain

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/achievement"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/progress"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain/uow"
)

type Container struct {
	progressSummaryRepo       progress.SummaryRepo
	topicStatRepo             progress.TopicStatRepo
	achievementDefinitionRepo achievement.DefinitionRepo
	candidateAchievementRepo  achievement.CandidateAchievementRepo
	uowFactory                uow.Factory
}

func NewContainer(
	progressSummaryRepo progress.SummaryRepo,
	topicStatRepo progress.TopicStatRepo,
	achievementDefinitionRepo achievement.DefinitionRepo,
	candidateAchievementRepo achievement.CandidateAchievementRepo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		progressSummaryRepo:       progressSummaryRepo,
		topicStatRepo:             topicStatRepo,
		achievementDefinitionRepo: achievementDefinitionRepo,
		candidateAchievementRepo:  candidateAchievementRepo,
		uowFactory:                uowFactory,
	}
}

func (c *Container) ProgressSummaryRepo() progress.SummaryRepo {
	return c.progressSummaryRepo
}

func (c *Container) TopicStatRepo() progress.TopicStatRepo {
	return c.topicStatRepo
}

func (c *Container) AchievementDefinitionRepo() achievement.DefinitionRepo {
	return c.achievementDefinitionRepo
}

func (c *Container) CandidateAchievementRepo() achievement.CandidateAchievementRepo {
	return c.candidateAchievementRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
