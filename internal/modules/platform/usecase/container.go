package usecase

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/alerterror/cleanuperrors"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/alerterror/geterror"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/alerterror/geterrorstats"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/alerterror/listerrors"
	tmcleanup "github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/cleanupresults"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/getqueuestats"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/listdlqtasks"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/listqueues"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/listschedules"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/listtaskresults"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/purgedlq"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/purgequeue"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/requeuefromdlq"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase/taskmill/triggerschedule"
)

type Container struct {
	listQueues      listqueues.UseCase
	getQueueStats   getqueuestats.UseCase
	listDLQTasks    listdlqtasks.UseCase
	listTaskResults listtaskresults.UseCase
	listSchedules   listschedules.UseCase
	requeueFromDLQ  requeuefromdlq.UseCase
	purgeQueue      purgequeue.UseCase
	purgeDLQ        purgedlq.UseCase
	cleanupResults  tmcleanup.UseCase
	triggerSchedule triggerschedule.UseCase

	listErrors    listerrors.UseCase
	getError      geterror.UseCase
	getErrorStats geterrorstats.UseCase
	cleanupErrors cleanuperrors.UseCase
}

func NewContainer(
	listQueues listqueues.UseCase,
	getQueueStats getqueuestats.UseCase,
	listDLQTasks listdlqtasks.UseCase,
	listTaskResults listtaskresults.UseCase,
	listSchedules listschedules.UseCase,
	requeueFromDLQ requeuefromdlq.UseCase,
	purgeQueue purgequeue.UseCase,
	purgeDLQ purgedlq.UseCase,
	cleanupResults tmcleanup.UseCase,
	triggerSchedule triggerschedule.UseCase,
	listErrors listerrors.UseCase,
	getError geterror.UseCase,
	getErrorStats geterrorstats.UseCase,
	cleanupErrors cleanuperrors.UseCase,
) *Container {
	return &Container{
		listQueues,
		getQueueStats,
		listDLQTasks,
		listTaskResults,
		listSchedules,
		requeueFromDLQ,
		purgeQueue,
		purgeDLQ,
		cleanupResults,
		triggerSchedule,
		listErrors,
		getError,
		getErrorStats,
		cleanupErrors,
	}
}

func (c *Container) ListQueues() listqueues.UseCase           { return c.listQueues }
func (c *Container) GetQueueStats() getqueuestats.UseCase     { return c.getQueueStats }
func (c *Container) ListDLQTasks() listdlqtasks.UseCase       { return c.listDLQTasks }
func (c *Container) ListTaskResults() listtaskresults.UseCase { return c.listTaskResults }
func (c *Container) ListSchedules() listschedules.UseCase     { return c.listSchedules }
func (c *Container) RequeueFromDLQ() requeuefromdlq.UseCase   { return c.requeueFromDLQ }
func (c *Container) PurgeQueue() purgequeue.UseCase           { return c.purgeQueue }
func (c *Container) PurgeDLQ() purgedlq.UseCase               { return c.purgeDLQ }
func (c *Container) CleanupResults() tmcleanup.UseCase        { return c.cleanupResults }
func (c *Container) TriggerSchedule() triggerschedule.UseCase { return c.triggerSchedule }

func (c *Container) ListErrors() listerrors.UseCase       { return c.listErrors }
func (c *Container) GetError() geterror.UseCase           { return c.getError }
func (c *Container) GetErrorStats() geterrorstats.UseCase { return c.getErrorStats }
func (c *Container) CleanupErrors() cleanuperrors.UseCase { return c.cleanupErrors }
