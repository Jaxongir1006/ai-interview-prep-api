package answer

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeAnswerNotFound         = "INTERVIEW_ANSWER_NOT_FOUND"
	CodeAnswerQuestionConflict = "INTERVIEW_ANSWER_QUESTION_CONFLICT"
)

type Answer struct {
	bun.BaseModel `bun:"table:interview_answers,alias:ia"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	SessionQuestionID int64 `json:"session_question_id"`

	AnswerText       string    `json:"answer_text"`
	TimeSpentSeconds int64     `json:"time_spent_seconds"`
	SubmittedAt      time.Time `json:"submitted_at"`
	CreatedAt        time.Time `json:"created_at"         bun:",nullzero"`
	UpdatedAt        time.Time `json:"updated_at"         bun:",nullzero"`
}

func (m *Answer) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
