package session

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeSessionNotFound   = "INTERVIEW_SESSION_NOT_FOUND"
	CodeSessionIDConflict = "INTERVIEW_SESSION_ID_CONFLICT"
)

const (
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusAbandoned  = "abandoned"
	StatusScoring    = "scoring"
)

const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
	DifficultyMixed  = "mixed"
)

type Session struct {
	bun.BaseModel `bun:"table:interview_sessions,alias:is"`

	ID string `json:"id" bun:"id,pk"`

	UserID string `json:"user_id"`

	Title           string   `json:"title"`
	TargetRole      string   `json:"target_role"`
	ExperienceLevel string   `json:"experience_level"`
	Difficulty      string   `json:"difficulty"`
	Status          string   `json:"status"`
	QuestionCount   int      `json:"question_count"`
	AnsweredCount   int      `json:"answered_count"`
	TotalScore      *float64 `json:"total_score"`

	TotalDurationSeconds int64      `json:"total_duration_seconds"`
	StartedAt            time.Time  `json:"started_at"`
	CompletedAt          *time.Time `json:"completed_at"`
	AbandonedAt          *time.Time `json:"abandoned_at"`
	CreatedAt            time.Time  `json:"created_at"             bun:",nullzero"`
	UpdatedAt            time.Time  `json:"updated_at"             bun:",nullzero"`
}

func (m *Session) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
