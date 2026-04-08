package profile

import "github.com/rise-and-shine/pkg/pg"

const (
	CodeCandidateProfileNotFound     = "CANDIDATE_PROFILE_NOT_FOUND"
	CodeCandidateProfileUserConflict = "CANDIDATE_PROFILE_USER_CONFLICT"
)

type CandidateProfile struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	UserID string `json:"user_id"`

	FullName *string `json:"full_name"`
	Bio      *string `json:"bio"`
	Location *string `json:"location"`

	TargetRole           *string `json:"target_role"`
	ExperienceLevel      *string `json:"experience_level"`
	InterviewGoalPerWeek int     `json:"interview_goal_per_week"`
}
