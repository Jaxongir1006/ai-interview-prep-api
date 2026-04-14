package http

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/forward"
)

type Controller struct {
	usecaseContainer *usecase.Container
	authPortal       auth.Portal
}

func NewController(
	usecaseContainer *usecase.Container,
	authPortal auth.Portal,
	httpServer *server.HTTPServer,
) *Controller {
	ctrl := &Controller{
		usecaseContainer: usecaseContainer,
		authPortal:       authPortal,
	}

	httpServer.RegisterRouter(ctrl.initRoutes)
	return ctrl
}

func (c *Controller) initRoutes(r fiber.Router) {
	v1 := r.Group("/api/v1/dashboard")
	v1Auth := v1.Group("", auth.NewAuthMiddleware(c.authPortal))

	v1Auth.Get("/overview", forward.ToUserAction(c.usecaseContainer.GetDashboardOverview()))
	v1Auth.Get("/stats", forward.ToUserAction(c.usecaseContainer.GetDashboardStats()))
	v1Auth.Get("/performance-trend", forward.ToUserAction(c.usecaseContainer.GetPerformanceTrend()))
	v1Auth.Get("/topics", forward.ToUserAction(c.usecaseContainer.GetDashboardTopics()))
	v1Auth.Get("/recent-activity", forward.ToUserAction(c.usecaseContainer.GetRecentActivity()))
	v1Auth.Get("/recommendations", forward.ToUserAction(c.usecaseContainer.GetDashboardRecommendations()))
}
