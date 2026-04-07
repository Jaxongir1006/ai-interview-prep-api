package http

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/usecase"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
)

type Controller struct {
	usecaseContainer *usecase.Container
	portalContainer  *portal.Container
}

func NewController(
	usecaseContainer *usecase.Container,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) *Controller {
	ctrl := &Controller{
		usecaseContainer: usecaseContainer,
		portalContainer:  portalContainer,
	}

	httpServer.RegisterRouter(ctrl.initRoutes)
	return ctrl
}

func (c *Controller) initRoutes(r fiber.Router) {
	v1 := r.Group("/api/v1/filevault", auth.NewAuthMiddleware(c.portalContainer.Auth()))

	v1.Post("/upload", c.Upload)
	v1.Get("/download", c.Download)
}
