package usecase

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/audit/usecase/actionlog/getactionlogs"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/audit/usecase/statuschangelog/getstatuschangelogs"
)

type Container struct {
	getActionLogs       getactionlogs.UseCase
	getStatusChangeLogs getstatuschangelogs.UseCase
}

func NewContainer(
	getActionLogs getactionlogs.UseCase,
	getStatusChangeLogs getstatuschangelogs.UseCase,
) *Container {
	return &Container{
		getActionLogs:       getActionLogs,
		getStatusChangeLogs: getStatusChangeLogs,
	}
}

func (c *Container) GetActionLogs() getactionlogs.UseCase { return c.getActionLogs }
func (c *Container) GetStatusChangeLogs() getstatuschangelogs.UseCase {
	return c.getStatusChangeLogs
}
