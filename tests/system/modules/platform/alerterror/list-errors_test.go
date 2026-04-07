//go:build system

package alerterror_test

import (
	"net/http"
	"testing"

	portalauth "github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/state/platform"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"
)

func TestListErrors_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionAlertView)
	platform.GivenErrors(t,
		map[string]any{"code": "INTERNAL_ERROR", "service": "svc-a", "operation": "op-1", "message": "error one"},
		map[string]any{"code": "VALIDATION_ERROR", "service": "svc-b", "operation": "op-2", "message": "error two"},
		map[string]any{"code": "INTERNAL_ERROR", "service": "svc-a", "operation": "op-3", "message": "error three"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/list-errors").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("count").Number().IsEqual(3)
	obj.Value("page_number").Number().IsEqual(1)
	obj.Value("page_size").Number().IsEqual(20)
	obj.Value("content").Array().Length().IsEqual(3)
}

func TestListErrors_FilterByService(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionAlertView)
	platform.GivenErrors(t,
		map[string]any{"service": "svc-a"},
		map[string]any{"service": "svc-a"},
		map[string]any{"service": "svc-b"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/list-errors").
		WithQuery("service", "svc-a").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("count").Number().IsEqual(2)
	obj.Value("content").Array().Length().IsEqual(2)
}

func TestListErrors_Search(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionAlertView)
	platform.GivenErrors(t,
		map[string]any{"message": "database connection failed"},
		map[string]any{"message": "auth token expired"},
		map[string]any{"code": "DATABASE_ERROR"},
	)

	// WHEN - search for "database"
	resp := trigger.UserAction(t).GET("/api/v1/platform/list-errors").
		WithQuery("search", "database").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN - matches message "database..." and code "DATABASE_ERROR"
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("count").Number().IsEqual(2)
}

func TestListErrors_AuthFailures(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing authorization header",
			setup:      func(_ *testing.T) string { return "" },
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name: "insufficient permissions",
			setup: func(t *testing.T) string {
				return auth.GivenAuthToken(t)
			},
			wantStatus: http.StatusForbidden,
			wantErr:    "FORBIDDEN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := tc.setup(t)

			// WHEN
			req := trigger.UserAction(t).GET("/api/v1/platform/list-errors")
			if token != "" {
				req = req.WithHeader("Authorization", "Bearer "+token)
			}
			resp := req.Expect()

			// THEN
			resp.Status(tc.wantStatus)
			resp.JSON().Object().Value("error").Object().
				Value("code").String().IsEqual(tc.wantErr)
			resp.JSON().Object().Value("error").Object().
				Value("message").String().NotContains("[untranslated]")
		})
	}
}
