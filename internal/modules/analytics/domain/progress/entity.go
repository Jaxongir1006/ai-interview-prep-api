package progress

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeProgressSummaryNotFound = "PROGRESS_SUMMARY_NOT_FOUND"
	CodeTopicStatNotFound       = "TOPIC_STAT_NOT_FOUND"
)

type Summary struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	CurrentStreak         int        `json:"current_streak"`
	LongestStreak         int        `json:"longest_streak"`
	TotalInterviewsTaken  int        `json:"total_interviews_taken"`
	TotalTimeSpentSeconds int64      `json:"total_time_spent_seconds"`
	AverageScore          float64    `json:"average_score"`
	LastInterviewAt       *time.Time `json:"last_interview_at"`
}

type TopicStat struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	TopicKey              string     `json:"topic_key"`
	Attempts              int        `json:"attempts"`
	TotalTimeSpentSeconds int64      `json:"total_time_spent_seconds"`
	AverageScore          float64    `json:"average_score"`
	BestScore             float64    `json:"best_score"`
	LastPracticedAt       *time.Time `json:"last_practiced_at"`
}
