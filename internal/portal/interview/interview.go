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

type CatalogTargetRole struct {
	Key          string
	Name         string
	Description  *string
	DisplayOrder int
}

type CatalogExperienceLevel struct {
	Key          string
	Name         string
	Description  *string
	DisplayOrder int
}

type CatalogTopic struct {
	Key            string
	Name           string
	Description    *string
	Category       *string
	TargetRoleKeys []string
	DisplayOrder   int
}

type GetOnboardingOptionsResponse struct {
	TargetRoles      []CatalogTargetRole
	ExperienceLevels []CatalogExperienceLevel
	Topics           []CatalogTopic
}

type ValidateOnboardingOptionsRequest struct {
	TargetRole      string
	ExperienceLevel string
	PreferredTopics []string
}

type ValidateOnboardingOptionsResponse struct {
	Valid                  bool
	UnknownTargetRole      bool
	UnknownExperienceLevel bool
	UnknownTopics          []string
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
	GetOnboardingOptions(ctx context.Context) (*GetOnboardingOptionsResponse, error)
	ValidateOnboardingOptions(
		ctx context.Context,
		req *ValidateOnboardingOptionsRequest,
	) (*ValidateOnboardingOptionsResponse, error)
	ListDashboardSessions(
		ctx context.Context,
		req *ListDashboardSessionsRequest,
	) (*ListDashboardSessionsResponse, error)
}
