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

func TestGetErrorStats_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionAlertView)
	platform.GivenErrors(t,
		map[string]any{"code": "INTERNAL_ERROR", "service": "svc-a", "operation": "op-1"},
		map[string]any{"code": "INTERNAL_ERROR", "service": "svc-a", "operation": "op-2"},
		map[string]any{"code": "VALIDATION_ERROR", "service": "svc-b", "operation": "op-1"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/get-error-stats").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("total_count").Number().IsEqual(3)
	obj.Value("by_service").Array().Length().IsEqual(2)
	obj.Value("by_operation").Array().Length().IsEqual(2)
	obj.Value("by_code").Array().Length().IsEqual(2)
}

func TestGetErrorStats_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/platform/get-error-stats")
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
