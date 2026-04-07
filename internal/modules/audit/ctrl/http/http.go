package http

import (
	"go-enterprise-blueprint/internal/modules/audit/usecase"
	"go-enterprise-blueprint/internal/portal"
	portalaudit "go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/internal/portal/auth"

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
	v1 := r.Group("/api/v1/audit")

	// All audit routes require authentication
	v1Auth := v1.Group("", auth.NewAuthMiddleware(c.authPortal))

	// Action logs
	actionLogRead := auth.RequirePermission(portalaudit.PermissionActionLogRead)
	v1Auth.Get("/get-action-logs", actionLogRead,
		forward.ToUserAction(c.usecaseContainer.GetActionLogs()))

	// Status change logs
	statusChangeLogRead := auth.RequirePermission(portalaudit.PermissionStatusChangeLogRead)
	v1Auth.Get("/get-status-change-logs", statusChangeLogRead,
		forward.ToUserAction(c.usecaseContainer.GetStatusChangeLogs()))
}
