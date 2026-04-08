package achievement

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeAchievementDefinitionNotFound = "ACHIEVEMENT_DEFINITION_NOT_FOUND"
	CodeCandidateAchievementNotFound  = "CANDIDATE_ACHIEVEMENT_NOT_FOUND"
	CodeAchievementCodeConflict       = "ACHIEVEMENT_CODE_CONFLICT"
)

type Definition struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	SortOrder   int     `json:"sort_order"`
}

type CandidateAchievement struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID                  string         `json:"user_id"`
	AchievementDefinitionID int64          `json:"achievement_definition_id"`
	AchievedAt              time.Time      `json:"achieved_at"`
	Metadata                map[string]any `json:"metadata"`
}
