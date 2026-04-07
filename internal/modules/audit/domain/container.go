package domain

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/actionlog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/statuschangelog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/uow"
)

// Container holds domain interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct {
	actionLogRepo       actionlog.Repo
	statusChangeLogRepo statuschangelog.Repo
	uowFactory          uow.Factory
}

func NewContainer(
	actionLogRepo actionlog.Repo,
	statusChangeLogRepo statuschangelog.Repo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		actionLogRepo,
		statusChangeLogRepo,
		uowFactory,
	}
}

func (c *Container) ActionLogRepo() actionlog.Repo {
	return c.actionLogRepo
}

func (c *Container) StatusChangeLogRepo() statuschangelog.Repo {
	return c.statusChangeLogRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
