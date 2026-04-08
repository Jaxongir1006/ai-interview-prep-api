package achievement

import "github.com/rise-and-shine/pkg/repogen"

type DefinitionFilter struct {
	ID       *int64
	Code     *string
	Category *string
	IDs      []int64
}

type CandidateAchievementFilter struct {
	ID                      *int64
	UserID                  *string
	AchievementDefinitionID *int64
	IDs                     []int64
}

type DefinitionRepo interface {
	repogen.Repo[Definition, DefinitionFilter]
}

type CandidateAchievementRepo interface {
	repogen.Repo[CandidateAchievement, CandidateAchievementFilter]
}
