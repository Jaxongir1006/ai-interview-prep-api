package cli

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/audit/usecase"
)

type Controller struct {
	usecaseContainer *usecase.Container
}

func NewController(usecaseContainer *usecase.Container) *Controller {
	return &Controller{
		usecaseContainer,
	}
}
