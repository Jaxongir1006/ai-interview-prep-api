package question

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeQuestionNotFound         = "INTERVIEW_QUESTION_NOT_FOUND"
	CodeQuestionPositionConflict = "INTERVIEW_QUESTION_POSITION_CONFLICT"
)

const (
	TypeTechnical    = "technical"
	TypeBehavioral   = "behavioral"
	TypeSystemDesign = "system_design"
	TypeCoding       = "coding"
)

const (
	SourceAI      = "ai"
	SourceManual  = "manual"
	SourceCatalog = "catalog"
)

type Question struct {
	bun.BaseModel `bun:"table:interview_questions,alias:iq"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	SessionID string `json:"session_id"`

	TopicKey         string         `json:"topic_key"`
	TopicName        string         `json:"topic_name"`
	Difficulty       string         `json:"difficulty"`
	QuestionType     string         `json:"question_type"`
	QuestionText     string         `json:"question_text"`
	ExpectedAnswer   *string        `json:"expected_answer"`
	EvaluationRubric map[string]any `json:"evaluation_rubric"  bun:",type:jsonb"`
	Source           string         `json:"source"`
	SourceQuestionID *string        `json:"source_question_id"`
	AIProvider       *string        `json:"ai_provider"`
	AIModel          *string        `json:"ai_model"`
	PromptVersion    *string        `json:"prompt_version"`
	Position         int            `json:"position"`
	CreatedAt        time.Time      `json:"created_at"         bun:",nullzero"`
	UpdatedAt        time.Time      `json:"updated_at"         bun:",nullzero"`
}

func (m *Question) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
