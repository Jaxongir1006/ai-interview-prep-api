package completeonboarding

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/rise-and-shine/pkg/val"
)

type Request struct {
	TargetRole      string   `json:"target_role"      validate:"required"`
	ExperienceLevel string   `json:"experience_level" validate:"required"`
	PreferredTopics []string `json:"preferred_topics" validate:"required,min=1,max=10,dive,required"`
}

type Response struct {
	Profile ProfileInfo `json:"profile"`
}

type ProfileInfo struct {
	ID                    int64     `json:"id"`
	UserID                string    `json:"user_id"`
	FullName              *string   `json:"full_name"`
	TargetRole            string    `json:"target_role"`
	ExperienceLevel       string    `json:"experience_level"`
	PreferredTopics       []string  `json:"preferred_topics"`
	OnboardingCompleted   bool      `json:"onboarding_completed"`
	OnboardingCompletedAt time.Time `json:"onboarding_completed_at"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, interviewPortal interview.Portal) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		interviewPortal: interviewPortal,
	}
}

type usecase struct {
	domainContainer *domain.Container
	interviewPortal interview.Portal
}

func (uc *usecase) OperationID() string { return "complete-onboarding" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	userCtx := auth.MustUserContext(ctx)
	if !userCtx.IsVerified {
		return nil, errx.New(
			"email is not verified",
			errx.WithType(errx.T_Validation),
			errx.WithCode(auth.CodeEmailNotVerified),
		)
	}

	if err := uc.validateInput(ctx, in); err != nil {
		return nil, errx.Wrap(err)
	}

	p, err := uc.domainContainer.ProfileRepo().Get(ctx, profile.Filter{
		UserID: &userCtx.UserID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, profile.CodeCandidateProfileNotFound)
	}

	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	now := time.Now()
	p.TargetRole = &in.TargetRole
	p.ExperienceLevel = &in.ExperienceLevel
	p.OnboardingCompleted = true
	p.OnboardingCompletedAt = &now

	p, err = uow.Profile().Update(ctx, p)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	err = uow.TopicPreference().DeleteByProfileID(ctx, p.ID)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	for i, topic := range in.PreferredTopics {
		_, err = uow.TopicPreference().Create(ctx, &topicpreference.TopicPreference{
			CandidateProfileID: p.ID,
			TopicKey:           topic,
			Priority:           i,
		})
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{
		Profile: ProfileInfo{
			ID:                    p.ID,
			UserID:                p.UserID,
			FullName:              p.FullName,
			TargetRole:            in.TargetRole,
			ExperienceLevel:       in.ExperienceLevel,
			PreferredTopics:       in.PreferredTopics,
			OnboardingCompleted:   p.OnboardingCompleted,
			OnboardingCompletedAt: *p.OnboardingCompletedAt,
		},
	}, nil
}

func (uc *usecase) validateInput(ctx context.Context, in *Request) error {
	seen := make(map[string]struct{}, len(in.PreferredTopics))
	for _, topic := range in.PreferredTopics {
		if _, ok := seen[topic]; ok {
			return errx.New(
				"preferred_topics must be unique",
				errx.WithType(errx.T_Validation),
				errx.WithCode(val.CodeValidationFailed),
			)
		}
		seen[topic] = struct{}{}
	}

	result, err := uc.interviewPortal.ValidateOnboardingOptions(ctx, &interview.ValidateOnboardingOptionsRequest{
		TargetRole:      in.TargetRole,
		ExperienceLevel: in.ExperienceLevel,
		PreferredTopics: in.PreferredTopics,
	})
	if err != nil {
		return errx.Wrap(err)
	}
	if result.Valid {
		return nil
	}

	if result.UnknownTargetRole {
		return errx.New(
			"target_role contains unknown or inactive value",
			errx.WithType(errx.T_Validation),
			errx.WithCode(val.CodeValidationFailed),
			errx.WithDetails(errx.D{
				"target_role": in.TargetRole,
			}),
		)
	}

	if result.UnknownExperienceLevel {
		return errx.New(
			"experience_level contains unknown or inactive value",
			errx.WithType(errx.T_Validation),
			errx.WithCode(val.CodeValidationFailed),
			errx.WithDetails(errx.D{
				"experience_level": in.ExperienceLevel,
			}),
		)
	}

	if len(result.UnknownTopics) > 0 {
		return errx.New(
			"preferred_topics contains unknown or inactive value",
			errx.WithType(errx.T_Validation),
			errx.WithCode(val.CodeValidationFailed),
			errx.WithDetails(errx.D{
				"topics": result.UnknownTopics,
			}),
		)
	}

	return nil
}
