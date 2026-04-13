package embassy

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/sorter"
)

func (e *embassy) CreateInitialProfile(
	ctx context.Context,
	req *candidateportal.CreateInitialProfileRequest,
) (*candidateportal.Profile, error) {
	var owned bool
	uow, err := e.domainContainer.UOWFactory().NewBorrowed(ctx)
	if err != nil {
		uow, err = e.domainContainer.UOWFactory().NewUOW(ctx)
		if err != nil {
			return nil, errx.Wrap(err)
		}
		owned = true
		defer uow.DiscardUnapplied()
	}

	p, err := uow.Profile().Create(ctx, &profile.CandidateProfile{
		UserID:               req.UserID,
		FullName:             req.FullName,
		InterviewGoalPerWeek: 3,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if owned {
		err = uow.ApplyChanges()
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	return toPortalProfile(p), nil
}

func (e *embassy) GetProfileByUserID(ctx context.Context, userID string) (*candidateportal.Profile, error) {
	p, err := e.domainContainer.ProfileRepo().Get(ctx, profile.Filter{
		UserID: &userID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return toPortalProfile(p), nil
}

func (e *embassy) ListTopicPreferencesByProfileID(
	ctx context.Context,
	profileID int64,
) ([]candidateportal.TopicPreference, error) {
	ps, err := e.domainContainer.TopicPreferenceRepo().List(ctx, topicpreference.Filter{
		CandidateProfileID: &profileID,
		SortOpts: sorter.SortOpts{
			{F: "priority", D: sorter.Asc},
			{F: "id", D: sorter.Asc},
		},
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	out := make([]candidateportal.TopicPreference, len(ps))
	for i := range ps {
		out[i] = candidateportal.TopicPreference{
			ID:        ps[i].ID,
			TopicKey:  ps[i].TopicKey,
			Priority:  ps[i].Priority,
			CreatedAt: ps[i].CreatedAt,
		}
	}

	return out, nil
}

func toPortalProfile(p *profile.CandidateProfile) *candidateportal.Profile {
	return &candidateportal.Profile{
		ID:                    p.ID,
		UserID:                p.UserID,
		FullName:              p.FullName,
		Bio:                   p.Bio,
		Location:              p.Location,
		TargetRole:            p.TargetRole,
		ExperienceLevel:       p.ExperienceLevel,
		InterviewGoalPerWeek:  p.InterviewGoalPerWeek,
		OnboardingCompleted:   p.OnboardingCompleted,
		OnboardingCompletedAt: p.OnboardingCompletedAt,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}
