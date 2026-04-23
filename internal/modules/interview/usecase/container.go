package usecase

import "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/usecase/catalog/getonboardingoptions"

type Container struct {
	getOnboardingOptions getonboardingoptions.UseCase
}

func NewContainer(
	getOnboardingOptions getonboardingoptions.UseCase,
) *Container {
	return &Container{
		getOnboardingOptions: getOnboardingOptions,
	}
}

func (c *Container) GetOnboardingOptions() getonboardingoptions.UseCase {
	return c.getOnboardingOptions
}
