package http

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/usecase"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/forward"
)

type Controller struct {
	usecaseContainer *usecase.Container
	portalContainer  *portal.Container
	authPortal       auth.Portal
}

func NewController(
	usecaseContainer *usecase.Container,
	portalContainer *portal.Container,
	authPortal auth.Portal,
	httpServer *server.HTTPServer,
) *Controller {
	ctrl := &Controller{
		usecaseContainer: usecaseContainer,
		portalContainer:  portalContainer,
		authPortal:       authPortal,
	}

	httpServer.RegisterRouter(ctrl.initRoutes)
	return ctrl
}

func (c *Controller) initRoutes(r fiber.Router) {
	v1 := r.Group("/api/v1/platform")

	// All platform routes require authentication
	v1Auth := v1.Group("", auth.NewAuthMiddleware(c.authPortal))

	// Taskmill console — read operations
	taskmillView := auth.RequirePermission(auth.PermissionTaskmillView)
	v1Auth.Get("/list-queues", taskmillView, forward.ToUserAction(c.usecaseContainer.ListQueues()))
	v1Auth.Get("/get-queue-stats", taskmillView, forward.ToUserAction(c.usecaseContainer.GetQueueStats()))
	v1Auth.Get("/list-dlq-tasks", taskmillView, forward.ToUserAction(c.usecaseContainer.ListDLQTasks()))
	v1Auth.Get("/list-task-results", taskmillView, forward.ToUserAction(c.usecaseContainer.ListTaskResults()))
	v1Auth.Get("/list-schedules", taskmillView, forward.ToUserAction(c.usecaseContainer.ListSchedules()))

	// Taskmill console — write operations
	taskmillManage := auth.RequirePermission(auth.PermissionTaskmillManage)
	v1Auth.Post("/requeue-from-dlq", taskmillManage, forward.ToUserAction(c.usecaseContainer.RequeueFromDLQ()))
	v1Auth.Post("/purge-queue", taskmillManage, forward.ToUserAction(c.usecaseContainer.PurgeQueue()))
	v1Auth.Post("/purge-dlq", taskmillManage, forward.ToUserAction(c.usecaseContainer.PurgeDLQ()))
	v1Auth.Post("/cleanup-results", taskmillManage, forward.ToUserAction(c.usecaseContainer.CleanupResults()))
	v1Auth.Post("/trigger-schedule", taskmillManage, forward.ToUserAction(c.usecaseContainer.TriggerSchedule()))

	// Alert errors — read operations
	alertView := auth.RequirePermission(auth.PermissionAlertView)
	v1Auth.Get("/list-errors", alertView, forward.ToUserAction(c.usecaseContainer.ListErrors()))
	v1Auth.Get("/get-error", alertView, forward.ToUserAction(c.usecaseContainer.GetError()))
	v1Auth.Get("/get-error-stats", alertView, forward.ToUserAction(c.usecaseContainer.GetErrorStats()))

	// Alert errors — write operations
	alertManage := auth.RequirePermission(auth.PermissionAlertManage)
	v1Auth.Post("/cleanup-errors", alertManage, forward.ToUserAction(c.usecaseContainer.CleanupErrors()))
}
