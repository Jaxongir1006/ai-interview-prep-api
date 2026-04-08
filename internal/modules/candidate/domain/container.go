package domain

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/uow"
)

type Container struct {
	profileRepo         profile.Repo
	topicPreferenceRepo topicpreference.Repo
	uowFactory          uow.Factory
}

func NewContainer(
	profileRepo profile.Repo,
	topicPreferenceRepo topicpreference.Repo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		profileRepo:         profileRepo,
		topicPreferenceRepo: topicPreferenceRepo,
		uowFactory:          uowFactory,
	}
}

func (c *Container) ProfileRepo() profile.Repo {
	return c.profileRepo
}

func (c *Container) TopicPreferenceRepo() topicpreference.Repo {
	return c.topicPreferenceRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
