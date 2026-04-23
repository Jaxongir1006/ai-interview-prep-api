package baseserver_test

import (
	"net/http"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/baseserver"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/stretchr/testify/require"
)

func TestCORSMiddlewareAllowsConfiguredLocalhostWildcardPorts(t *testing.T) {
	app := newTestApp(baseserver.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:*",
			"http://127.0.0.1:*",
		},
		AllowMethods:     []string{http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	req, err := http.NewRequest(http.MethodOptions, "/api/v1/auth/register", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:5174")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "content-type,authorization")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Equal(t, "http://localhost:5174", resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
	require.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), http.MethodPost)
	require.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Content-Type")
	require.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Authorization")
	require.Equal(t, "86400", resp.Header.Get("Access-Control-Max-Age"))
}

func TestCORSMiddlewareRejectsNonLocalhostWildcardOrigins(t *testing.T) {
	app := newTestApp(baseserver.CORSConfig{
		AllowOrigins:     []string{"http://localhost:*"},
		AllowMethods:     []string{http.MethodPost, http.MethodOptions},
		AllowCredentials: true,
	})

	req, err := http.NewRequest(http.MethodOptions, "/api/v1/auth/register", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://example.com:5174")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	require.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
}

func newTestApp(corsCfg baseserver.CORSConfig) *fiber.App {
	corsCfg.Enabled = true
	srv := baseserver.New(server.Config{}, corsCfg)
	app := srv.GetApp()
	app.Post("/api/v1/auth/register", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	return app
}
