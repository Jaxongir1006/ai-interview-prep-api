package domain

import (
	"go-enterprise-blueprint/internal/modules/platform/domain/alerterror"
	"go-enterprise-blueprint/internal/modules/platform/domain/uow"

	"github.com/rise-and-shine/pkg/taskmill/console"
)

type Container struct {
	console        console.Console
	alertErrorRepo alerterror.Repo
	uowFactory     uow.Factory
}

func NewContainer(console console.Console, alertErrorRepo alerterror.Repo, uowFactory uow.Factory) *Container {
	return &Container{
		console:        console,
		alertErrorRepo: alertErrorRepo,
		uowFactory:     uowFactory,
	}
}

func (c *Container) Console() console.Console {
	return c.console
}

func (c *Container) AlertErrorRepo() alerterror.Repo {
	return c.alertErrorRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
