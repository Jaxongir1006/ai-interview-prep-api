package review

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeReviewNotFound         = "INTERVIEW_REVIEW_NOT_FOUND"
	CodeReviewQuestionConflict = "INTERVIEW_REVIEW_QUESTION_CONFLICT"
)

const (
	ReviewerTypeAI     = "ai"
	ReviewerTypeManual = "manual"
)

type Review struct {
	bun.BaseModel `bun:"table:interview_reviews,alias:ir"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	SessionQuestionID int64  `json:"session_question_id"`
	AnswerID          *int64 `json:"answer_id"`

	ReviewerType    string         `json:"reviewer_type"`
	Score           *float64       `json:"score"`
	CorrectnessRate *float64       `json:"correctness_rate"`
	Feedback        *string        `json:"feedback"`
	RubricScores    map[string]any `json:"rubric_scores"    bun:",type:jsonb"`
	Strengths       map[string]any `json:"strengths"        bun:",type:jsonb"`
	Improvements    map[string]any `json:"improvements"     bun:",type:jsonb"`
	AIProvider      *string        `json:"ai_provider"`
	AIModel         *string        `json:"ai_model"`
	PromptVersion   *string        `json:"prompt_version"`
	Metadata        map[string]any `json:"metadata"         bun:",type:jsonb"`
	ReviewedAt      time.Time      `json:"reviewed_at"`
	CreatedAt       time.Time      `json:"created_at"       bun:",nullzero"`
	UpdatedAt       time.Time      `json:"updated_at"       bun:",nullzero"`
}

func (m *Review) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
