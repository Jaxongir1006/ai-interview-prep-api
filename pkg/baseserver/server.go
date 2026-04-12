// Package baseserver provides a constructor function to create a base HTTP server with standard middlewares.
package baseserver

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/middleware"
)

// CORSConfig defines the browser cross-origin policy for the HTTP API.
type CORSConfig struct {
	Enabled bool `yaml:"enabled"`

	AllowOrigins  []string `yaml:"allow_origins"`
	AllowMethods  []string `yaml:"allow_methods"`
	AllowHeaders  []string `yaml:"allow_headers"`
	ExposeHeaders []string `yaml:"expose_headers"`

	AllowCredentials bool `yaml:"allow_credentials"`
	MaxAge           int  `yaml:"max_age"`
}

func New(
	cfg server.Config,
	corsCfg CORSConfig,
) *server.HTTPServer {
	middlewares := []server.Middleware{
		middleware.NewRecoveryMW(cfg.HideErrorDetails),
		middleware.NewTracingMW(),
		middleware.NewTimeoutMW(cfg.HandleTimeout),
		middleware.NewAlertingMW(),
		middleware.NewLoggerMW(cfg.HideErrorDetails),
		middleware.NewErrorHandlerMW(cfg.HideErrorDetails),
	}

	if corsCfg.Enabled {
		middlewares = append(middlewares, newCORSMiddleware(corsCfg))
	}

	return server.NewHTTPServer(cfg, middlewares)
}

func newCORSMiddleware(cfg CORSConfig) server.Middleware {
	return server.Middleware{
		Priority: 950,
		Handler: fibercors.New(fibercors.Config{
			AllowOrigins:     joinOrDefault(cfg.AllowOrigins, "*"),
			AllowMethods:     joinOrDefault(cfg.AllowMethods, strings.Join(defaultAllowMethods(), ",")),
			AllowHeaders:     strings.Join(cfg.AllowHeaders, ","),
			AllowCredentials: cfg.AllowCredentials,
			ExposeHeaders:    strings.Join(cfg.ExposeHeaders, ","),
			MaxAge:           cfg.MaxAge,
		}),
	}
}

func defaultAllowMethods() []string {
	return []string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodOptions,
	}
}

func joinOrDefault(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}

	return strings.Join(values, ",")
}
