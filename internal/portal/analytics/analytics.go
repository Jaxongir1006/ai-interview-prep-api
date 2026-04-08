package analytics

import (
	"context"
	"time"
)

const (
	ModuleName                  = "analytics"
	CodeProgressSummaryNotFound = "PROGRESS_SUMMARY_NOT_FOUND"
)

type ProgressSummary struct {
	ID                    int64
	UserID                string
	CurrentStreak         int
	LongestStreak         int
	TotalInterviewsTaken  int
	TotalTimeSpentSeconds int64
	AverageScore          float64
	LastInterviewAt       *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type Portal interface {
	GetProgressSummaryByUserID(ctx context.Context, userID string) (*ProgressSummary, error)
}
