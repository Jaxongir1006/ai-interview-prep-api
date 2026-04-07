package cli

import (
	"go-enterprise-blueprint/internal/modules/platform/usecase"
)

type Controller struct {
	usecaseContainer *usecase.Container
}

func NewController(usecaseContainer *usecase.Container) *Controller {
	return &Controller{
		usecaseContainer,
	}
}
