// Package baseserver provides a constructor function to create a base HTTP server with standard middlewares.
package baseserver

import (
	"net/url"
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
	allowOrigins, allowOriginFunc := corsAllowOrigins(cfg.AllowOrigins)

	return server.Middleware{
		Priority: 950,
		Handler: fibercors.New(fibercors.Config{
			AllowOrigins:     corsAllowOriginsValue(allowOrigins, allowOriginFunc),
			AllowOriginsFunc: allowOriginFunc,
			AllowMethods:     joinOrDefault(cfg.AllowMethods, strings.Join(defaultAllowMethods(), ",")),
			AllowHeaders:     strings.Join(cfg.AllowHeaders, ","),
			AllowCredentials: cfg.AllowCredentials,
			ExposeHeaders:    strings.Join(cfg.ExposeHeaders, ","),
			MaxAge:           cfg.MaxAge,
		}),
	}
}

func corsAllowOriginsValue(origins []string, allowOriginFunc func(string) bool) string {
	if len(origins) == 0 && allowOriginFunc != nil {
		return ""
	}

	return joinOrDefault(origins, "*")
}

func corsAllowOrigins(origins []string) ([]string, func(string) bool) {
	staticOrigins := make([]string, 0, len(origins))
	localhostWildcards := make([]string, 0)

	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if isLocalhostWildcardOrigin(trimmed) {
			localhostWildcards = append(localhostWildcards, strings.TrimSuffix(trimmed, ":*"))
			continue
		}

		staticOrigins = append(staticOrigins, trimmed)
	}

	if len(localhostWildcards) == 0 {
		return staticOrigins, nil
	}

	return staticOrigins, func(origin string) bool {
		u, err := url.Parse(strings.ToLower(origin))
		if err != nil {
			return false
		}

		if u.Scheme == "" || u.Hostname() == "" || u.Port() == "" {
			return false
		}

		normalized := u.Scheme + "://" + u.Hostname()
		for _, allowed := range localhostWildcards {
			if normalized == allowed {
				return true
			}
		}

		return false
	}
}

func isLocalhostWildcardOrigin(origin string) bool {
	origin = strings.ToLower(strings.TrimSpace(origin))
	if !strings.HasSuffix(origin, ":*") {
		return false
	}

	return strings.TrimSuffix(origin, ":*") == "http://localhost" ||
		strings.TrimSuffix(origin, ":*") == "https://localhost" ||
		strings.TrimSuffix(origin, ":*") == "http://127.0.0.1" ||
		strings.TrimSuffix(origin, ":*") == "https://127.0.0.1"
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
