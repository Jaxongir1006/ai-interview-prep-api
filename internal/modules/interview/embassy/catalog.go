package embassy

import (
	"context"
	"slices"

	interviewportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/code19m/errx"
)

func (e *embassy) GetOnboardingOptions(ctx context.Context) (*interviewportal.GetOnboardingOptionsResponse, error) {
	roles, err := e.domainContainer.CatalogRepo().ListActiveTargetRoles(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	levels, err := e.domainContainer.CatalogRepo().ListActiveExperienceLevels(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	topics, err := e.domainContainer.CatalogRepo().ListActiveTopics(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	mappings, err := e.domainContainer.CatalogRepo().ListActiveTargetRoleTopics(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	roleKeysByID := make(map[int64]string, len(roles))
	for i := range roles {
		roleKeysByID[roles[i].ID] = roles[i].Key
	}

	targetRoleKeysByTopicID := map[int64][]string{}
	for i := range mappings {
		roleKey, ok := roleKeysByID[mappings[i].TargetRoleID]
		if !ok {
			continue
		}
		targetRoleKeysByTopicID[mappings[i].TopicID] = append(targetRoleKeysByTopicID[mappings[i].TopicID], roleKey)
	}

	out := &interviewportal.GetOnboardingOptionsResponse{
		TargetRoles:      make([]interviewportal.CatalogTargetRole, 0, len(roles)),
		ExperienceLevels: make([]interviewportal.CatalogExperienceLevel, 0, len(levels)),
		Topics:           make([]interviewportal.CatalogTopic, 0, len(topics)),
	}

	for i := range roles {
		out.TargetRoles = append(out.TargetRoles, interviewportal.CatalogTargetRole{
			Key:          roles[i].Key,
			Name:         roles[i].Name,
			Description:  roles[i].Description,
			DisplayOrder: roles[i].DisplayOrder,
		})
	}
	for i := range levels {
		out.ExperienceLevels = append(out.ExperienceLevels, interviewportal.CatalogExperienceLevel{
			Key:          levels[i].Key,
			Name:         levels[i].Name,
			Description:  levels[i].Description,
			DisplayOrder: levels[i].DisplayOrder,
		})
	}
	for i := range topics {
		out.Topics = append(out.Topics, interviewportal.CatalogTopic{
			Key:            topics[i].Key,
			Name:           topics[i].Name,
			Description:    topics[i].Description,
			Category:       topics[i].Category,
			TargetRoleKeys: targetRoleKeysByTopicID[topics[i].ID],
			DisplayOrder:   topics[i].DisplayOrder,
		})
	}

	return out, nil
}

func (e *embassy) ValidateOnboardingOptions(
	ctx context.Context,
	req *interviewportal.ValidateOnboardingOptionsRequest,
) (*interviewportal.ValidateOnboardingOptionsResponse, error) {
	options, err := e.GetOnboardingOptions(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	roleKeys := make([]string, 0, len(options.TargetRoles))
	for i := range options.TargetRoles {
		roleKeys = append(roleKeys, options.TargetRoles[i].Key)
	}

	levelKeys := make([]string, 0, len(options.ExperienceLevels))
	for i := range options.ExperienceLevels {
		levelKeys = append(levelKeys, options.ExperienceLevels[i].Key)
	}

	topicKeys := make([]string, 0, len(options.Topics))
	for i := range options.Topics {
		topicKeys = append(topicKeys, options.Topics[i].Key)
	}

	out := &interviewportal.ValidateOnboardingOptionsResponse{
		UnknownTargetRole:      !slices.Contains(roleKeys, req.TargetRole),
		UnknownExperienceLevel: !slices.Contains(levelKeys, req.ExperienceLevel),
	}
	for _, topic := range req.PreferredTopics {
		if !slices.Contains(topicKeys, topic) {
			out.UnknownTopics = append(out.UnknownTopics, topic)
		}
	}
	out.Valid = !out.UnknownTargetRole && !out.UnknownExperienceLevel && len(out.UnknownTopics) == 0

	return out, nil
}
