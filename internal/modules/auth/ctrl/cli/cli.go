// Package cli provides container of cobra CLI commands for auth module.
package cli

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase"
)

type Controller struct {
	usecaseContainer *usecase.Container
}

func NewController(usecaseContainer *usecase.Container) *Controller {
	return &Controller{
		usecaseContainer,
	}
}
