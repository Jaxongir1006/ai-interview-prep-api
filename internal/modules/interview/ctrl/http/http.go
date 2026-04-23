package http

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/usecase"
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
	v1 := r.Group("/api/v1/interview")
	v1Auth := v1.Group("", auth.NewAuthMiddleware(c.authPortal))

	v1Auth.Get("/get-onboarding-options", forward.ToUserAction(c.usecaseContainer.GetOnboardingOptions()))
}
