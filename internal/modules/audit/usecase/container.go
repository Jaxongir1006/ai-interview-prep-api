package usecase

import (
	"go-enterprise-blueprint/internal/modules/audit/usecase/actionlog/getactionlogs"
	"go-enterprise-blueprint/internal/modules/audit/usecase/statuschangelog/getstatuschangelogs"
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
