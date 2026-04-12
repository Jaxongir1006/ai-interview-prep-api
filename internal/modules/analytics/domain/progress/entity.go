package progress

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeProgressSummaryNotFound = "PROGRESS_SUMMARY_NOT_FOUND"
	CodeTopicStatNotFound       = "TOPIC_STAT_NOT_FOUND"
)

type Summary struct {
	bun.BaseModel `bun:"table:candidate_progress_summaries,alias:cps"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	CurrentStreak         int        `json:"current_streak"`
	LongestStreak         int        `json:"longest_streak"`
	TotalInterviewsTaken  int        `json:"total_interviews_taken"`
	TotalTimeSpentSeconds int64      `json:"total_time_spent_seconds"`
	AverageScore          float64    `json:"average_score"`
	LastInterviewAt       *time.Time `json:"last_interview_at"`

	CreatedAt time.Time `bun:",nullzero" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero" json:"updated_at"`
}

type TopicStat struct {
	bun.BaseModel `bun:"table:candidate_topic_stats,alias:cts"`

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	TopicKey              string     `json:"topic_key"`
	Attempts              int        `json:"attempts"`
	TotalTimeSpentSeconds int64      `json:"total_time_spent_seconds"`
	AverageScore          float64    `json:"average_score"`
	BestScore             float64    `json:"best_score"`
	LastPracticedAt       *time.Time `json:"last_practiced_at"`

	CreatedAt time.Time `bun:",nullzero" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero" json:"updated_at"`
}

func (m *Summary) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

func (m *TopicStat) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
