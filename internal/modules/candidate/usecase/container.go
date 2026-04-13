package usecase

import "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/usecase/profile/completeonboarding"

type Container struct {
	completeOnboarding completeonboarding.UseCase
}

func NewContainer(
	completeOnboarding completeonboarding.UseCase,
) *Container {
	return &Container{
		completeOnboarding: completeOnboarding,
	}
}

func (c *Container) CompleteOnboarding() completeonboarding.UseCase {
	return c.completeOnboarding
}
