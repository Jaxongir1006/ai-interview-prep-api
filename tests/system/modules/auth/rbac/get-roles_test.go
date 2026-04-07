//go:build system

package rbac_test

import (
	"net/http"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"
)

func TestGetRoles_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	auth.GivenRoles(t,
		map[string]any{"name": "admin"},
		map[string]any{"name": "editor"},
		map[string]any{"name": "viewer"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-roles").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("page_number").Number().IsEqual(1)
	obj.Value("page_size").Number().IsEqual(20)
	obj.Value("count").Number().IsEqual(3)
	obj.Value("content").Array().Length().IsEqual(3)
}

func TestGetRoles_Empty(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-roles").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("count").Number().IsEqual(0)
	obj.Value("content").Array().Length().IsEqual(0)
}

func TestGetRoles_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/auth/get-roles")
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
