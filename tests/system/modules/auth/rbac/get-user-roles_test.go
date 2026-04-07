//go:build system

package rbac_test

import (
	"net/http"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"
)

func TestGetUserRoles_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	roles := auth.GivenRoles(t,
		map[string]any{"name": "admin"},
		map[string]any{"name": "editor"},
	)
	auth.GivenUserRoles(t,
		map[string]any{"user_id": u.ID, "role_id": roles[0].ID},
		map[string]any{"user_id": u.ID, "role_id": roles[1].ID},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-user-roles").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("user_id", u.ID).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().IsEqual(2)
}

func TestGetUserRoles_UserNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-user-roles").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("user_id", "00000000-0000-0000-0000-000000000000").
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestGetUserRoles_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/auth/get-user-roles")
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
