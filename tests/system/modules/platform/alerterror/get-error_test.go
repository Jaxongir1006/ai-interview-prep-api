//go:build system

package alerterror_test

import (
	"net/http"
	"testing"

	portalauth "go-enterprise-blueprint/internal/portal/auth"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/state/platform"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/google/uuid"
)

func TestGetError_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionAlertView)
	ids := platform.GivenErrors(t, map[string]any{
		"code":      "INTERNAL_ERROR",
		"service":   "my-svc",
		"operation": "test-op",
		"message":   "something failed",
		"details":   map[string]string{"trace_id": "abc-123"},
	})

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/get-error").
		WithQuery("id", ids[0]).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("id").String().IsEqual(ids[0])
	obj.Value("code").String().IsEqual("INTERNAL_ERROR")
	obj.Value("service").String().IsEqual("my-svc")
	obj.Value("operation").String().IsEqual("test-op")
	obj.Value("message").String().IsEqual("something failed")
	obj.Value("details").Object().Value("trace_id").String().IsEqual("abc-123")
	obj.ContainsKey("created_at")
	obj.Value("alerted").Boolean().IsFalse()
}

func TestGetError_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionAlertView)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/get-error").
		WithQuery("id", uuid.NewString()).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("ERROR_NOT_FOUND")
}

func TestGetError_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/platform/get-error").
				WithQuery("id", uuid.NewString())
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
