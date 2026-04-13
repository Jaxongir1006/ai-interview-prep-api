package getme

import (
	"context"
	"fmt"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/oauthaccount"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	analyticsportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/analytics"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"
	filevaultportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/filevault"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const avatarAssocType = "avatar"

type Request struct{}

type Response struct {
	User            UserInfo             `json:"user"`
	Profile         *ProfileInfo         `json:"profile"`
	ProgressSummary *ProgressSummaryInfo `json:"progress_summary"`
	Avatar          *AvatarInfo          `json:"avatar"`
}

type UserInfo struct {
	ID             string   `json:"id"`
	Username       *string  `json:"username"`
	Email          *string  `json:"email"`
	PhoneNumber    *string  `json:"phone_number"`
	IsVerified     bool     `json:"is_verified"`
	IsActive       bool     `json:"is_active"`
	LastLoginAt    any      `json:"last_login_at"`
	LastActiveAt   any      `json:"last_active_at"`
	CreatedAt      any      `json:"created_at"`
	UpdatedAt      any      `json:"updated_at"`
	OAuthProviders []string `json:"oauth_providers"`
}

type ProfileInfo struct {
	ID                    int64    `json:"id"`
	UserID                string   `json:"user_id"`
	FullName              *string  `json:"full_name"`
	Bio                   *string  `json:"bio"`
	Location              *string  `json:"location"`
	TargetRole            *string  `json:"target_role"`
	ExperienceLevel       *string  `json:"experience_level"`
	InterviewGoalPerWeek  int      `json:"interview_goal_per_week"`
	PreferredTopics       []string `json:"preferred_topics"`
	OnboardingCompleted   bool     `json:"onboarding_completed"`
	OnboardingCompletedAt any      `json:"onboarding_completed_at"`
	CreatedAt             any      `json:"created_at"`
	UpdatedAt             any      `json:"updated_at"`
}

type ProgressSummaryInfo struct {
	CurrentStreak         int     `json:"current_streak"`
	LongestStreak         int     `json:"longest_streak"`
	TotalInterviewsTaken  int     `json:"total_interviews_taken"`
	TotalTimeSpentSeconds int64   `json:"total_time_spent_seconds"`
	AverageScore          float64 `json:"average_score"`
	LastInterviewAt       any     `json:"last_interview_at"`
}

type AvatarInfo struct {
	FileID           string `json:"file_id"`
	OriginalFilename string `json:"original_filename"`
	MimeType         string `json:"mime_type"`
	SizeBytes        int64  `json:"size_bytes"`
	DownloadURL      string `json:"download_url"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		portalContainer: portalContainer,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "get-me" }

func (uc *usecase) Execute(ctx context.Context, _ *Request) (*Response, error) {
	userCtx := auth.MustUserContext(ctx)

	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		ID: &userCtx.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	providers, err := uc.listOAuthProviders(ctx, userCtx.UserID)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	resp := &Response{
		User: UserInfo{
			ID:             u.ID,
			Username:       u.Username,
			Email:          u.Email,
			PhoneNumber:    u.PhoneNumber,
			IsVerified:     u.IsVerified,
			IsActive:       u.IsActive,
			LastLoginAt:    u.LastLoginAt,
			LastActiveAt:   u.LastActiveAt,
			CreatedAt:      u.CreatedAt,
			UpdatedAt:      u.UpdatedAt,
			OAuthProviders: providers,
		},
	}

	p, err := uc.portalContainer.Candidate().GetProfileByUserID(ctx, userCtx.UserID)
	if err != nil && !errx.IsCodeIn(err, candidateportal.CodeProfileNotFound) {
		return nil, errx.Wrap(err)
	}
	if err == nil {
		preferredTopics, prefErr := uc.portalContainer.Candidate().ListTopicPreferencesByProfileID(ctx, p.ID)
		if prefErr != nil {
			return nil, errx.Wrap(prefErr)
		}

		resp.Profile = &ProfileInfo{
			ID:                    p.ID,
			UserID:                p.UserID,
			FullName:              p.FullName,
			Bio:                   p.Bio,
			Location:              p.Location,
			TargetRole:            p.TargetRole,
			ExperienceLevel:       p.ExperienceLevel,
			InterviewGoalPerWeek:  p.InterviewGoalPerWeek,
			PreferredTopics:       toTopicKeys(preferredTopics),
			OnboardingCompleted:   p.OnboardingCompleted,
			OnboardingCompletedAt: p.OnboardingCompletedAt,
			CreatedAt:             p.CreatedAt,
			UpdatedAt:             p.UpdatedAt,
		}

		avatar, foundAvatar, avatarErr := uc.findAvatar(ctx, p.ID)
		if avatarErr != nil {
			return nil, errx.Wrap(avatarErr)
		}
		if foundAvatar {
			resp.Avatar = avatar
		}
	}

	summary, err := uc.portalContainer.Analytics().GetProgressSummaryByUserID(ctx, userCtx.UserID)
	if err != nil && !errx.IsCodeIn(err, analyticsportal.CodeProgressSummaryNotFound) {
		return nil, errx.Wrap(err)
	}
	if err == nil {
		resp.ProgressSummary = &ProgressSummaryInfo{
			CurrentStreak:         summary.CurrentStreak,
			LongestStreak:         summary.LongestStreak,
			TotalInterviewsTaken:  summary.TotalInterviewsTaken,
			TotalTimeSpentSeconds: summary.TotalTimeSpentSeconds,
			AverageScore:          summary.AverageScore,
			LastInterviewAt:       summary.LastInterviewAt,
		}
	}

	return resp, nil
}

func (uc *usecase) listOAuthProviders(ctx context.Context, userID string) ([]string, error) {
	accounts, err := uc.domainContainer.OAuthAccountRepo().List(ctx, oauthaccount.Filter{
		UserID: &userID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	providers := make([]string, len(accounts))
	for i := range accounts {
		providers[i] = accounts[i].Provider
	}

	return providers, nil
}

func (uc *usecase) findAvatar(ctx context.Context, profileID int64) (*AvatarInfo, bool, error) {
	files, err := uc.portalContainer.Filevault().ListByEntity(ctx, &filevaultportal.ListByEntityRequest{
		EntityType: candidateportal.EntityTypeProfile,
		EntityID:   profileID,
		AssocType:  strPtr(avatarAssocType),
	})
	if err != nil {
		return nil, false, errx.Wrap(err)
	}
	if len(files) == 0 {
		return nil, false, nil
	}

	f := files[0]
	return &AvatarInfo{
		FileID:           f.ID,
		OriginalFilename: f.OriginalName,
		MimeType:         f.ContentType,
		SizeBytes:        f.Size,
		DownloadURL:      fmt.Sprintf("/api/v1/filevault/download?id=%s", f.ID),
	}, true, nil
}

func toTopicKeys(items []candidateportal.TopicPreference) []string {
	out := make([]string, len(items))
	for i := range items {
		out[i] = items[i].TopicKey
	}
	return out
}

func strPtr(v string) *string { return &v }
