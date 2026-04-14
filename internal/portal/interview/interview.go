package interview

import (
	"context"
	"time"
)

const ModuleName = "interview"

type Topic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DashboardSession struct {
	ID                    string
	Title                 string
	Status                string
	Score                 *float64
	StartedAt             time.Time
	CompletedAt           *time.Time
	DurationSeconds       int64
	QuestionCount         int
	AnsweredCount         int
	Topics                []Topic
	ReviewedQuestionCount int
}

type ListDashboardSessionsRequest struct {
	UserID        string
	Statuses      []string
	StartedAtFrom *time.Time
	StartedAtTo   *time.Time
	TopicID       *string
	Limit         int
	Cursor        *string
}

type ListDashboardSessionsResponse struct {
	Items      []DashboardSession
	NextCursor *string
}

type Portal interface {
	ListDashboardSessions(
		ctx context.Context,
		req *ListDashboardSessionsRequest,
	) (*ListDashboardSessionsResponse, error)
}
