package candidate

import (
	"context"
	"time"
)

const (
	ModuleName          = "candidate"
	CodeProfileNotFound = "CANDIDATE_PROFILE_NOT_FOUND"

	EntityTypeProfile = "candidate_profile"
)

type Profile struct {
	ID                    int64
	UserID                string
	FullName              *string
	Bio                   *string
	Location              *string
	TargetRole            *string
	ExperienceLevel       *string
	InterviewGoalPerWeek  int
	OnboardingCompleted   bool
	OnboardingCompletedAt *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type TopicPreference struct {
	ID        int64
	TopicKey  string
	Priority  int
	CreatedAt time.Time
}

type CreateInitialProfileRequest struct {
	UserID   string
	FullName *string
}

type Portal interface {
	CreateInitialProfile(ctx context.Context, req *CreateInitialProfileRequest) (*Profile, error)
	GetProfileByUserID(ctx context.Context, userID string) (*Profile, error)
	ListTopicPreferencesByProfileID(ctx context.Context, profileID int64) ([]TopicPreference, error)
}
