package getonboardingoptions

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct{}

type Response struct {
	TargetRoles      []TargetRole      `json:"target_roles"`
	ExperienceLevels []ExperienceLevel `json:"experience_levels"`
	Topics           []Topic           `json:"topics"`
}

type TargetRole struct {
	Key          string  `json:"key"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	DisplayOrder int     `json:"display_order"`
}

type ExperienceLevel struct {
	Key          string  `json:"key"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	DisplayOrder int     `json:"display_order"`
}

type Topic struct {
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	Category       *string  `json:"category"`
	TargetRoleKeys []string `json:"target_role_keys"`
	DisplayOrder   int      `json:"display_order"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(interviewPortal interview.Portal) UseCase {
	return &usecase{interviewPortal: interviewPortal}
}

type usecase struct {
	interviewPortal interview.Portal
}

func (uc *usecase) OperationID() string { return "get-onboarding-options" }

func (uc *usecase) Execute(ctx context.Context, _ *Request) (*Response, error) {
	options, err := uc.interviewPortal.GetOnboardingOptions(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	out := &Response{
		TargetRoles:      make([]TargetRole, 0, len(options.TargetRoles)),
		ExperienceLevels: make([]ExperienceLevel, 0, len(options.ExperienceLevels)),
		Topics:           make([]Topic, 0, len(options.Topics)),
	}
	for i := range options.TargetRoles {
		out.TargetRoles = append(out.TargetRoles, TargetRole{
			Key:          options.TargetRoles[i].Key,
			Name:         options.TargetRoles[i].Name,
			Description:  options.TargetRoles[i].Description,
			DisplayOrder: options.TargetRoles[i].DisplayOrder,
		})
	}
	for i := range options.ExperienceLevels {
		out.ExperienceLevels = append(out.ExperienceLevels, ExperienceLevel{
			Key:          options.ExperienceLevels[i].Key,
			Name:         options.ExperienceLevels[i].Name,
			Description:  options.ExperienceLevels[i].Description,
			DisplayOrder: options.ExperienceLevels[i].DisplayOrder,
		})
	}
	for i := range options.Topics {
		out.Topics = append(out.Topics, Topic{
			Key:            options.Topics[i].Key,
			Name:           options.Topics[i].Name,
			Description:    options.Topics[i].Description,
			Category:       options.Topics[i].Category,
			TargetRoleKeys: options.Topics[i].TargetRoleKeys,
			DisplayOrder:   options.Topics[i].DisplayOrder,
		})
	}

	return out, nil
}
