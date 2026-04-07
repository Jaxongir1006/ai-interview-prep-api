//go:build system

package user_test

import (
	"net/http"
	"testing"
	"time"

	portalauth "go-enterprise-blueprint/internal/portal/auth"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"
)

func TestGetAuthStats_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionUserRead)

	// Create additional users (GivenAuthToken already created 1)
	users := auth.GivenUsers(t,
		map[string]any{"username": "user-two"},
		map[string]any{"username": "user-three"},
	)

	// Create roles
	auth.GivenRoles(t,
		map[string]any{"name": "admin"},
		map[string]any{"name": "editor"},
	)

	// Create additional sessions (GivenAuthToken already created 1 active session)
	auth.GivenSessions(t,
		map[string]any{"user_id": users[0].ID},
		map[string]any{"user_id": users[1].ID},
	)

	// Create an expired session (should not be counted)
	auth.GivenSessions(t, map[string]any{
		"user_id":                  users[0].ID,
		"refresh_token_expires_at": time.Now().Add(-1 * time.Hour),
	})

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-auth-stats").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("total_users").Number().IsEqual(3)     // 1 from GivenAuthToken + 2 additional
	obj.Value("total_roles").Number().IsEqual(2)     // admin + editor
	obj.Value("active_sessions").Number().IsEqual(3) // 1 from GivenAuthToken + 2 additional (expired not counted)
}

func TestGetAuthStats_Empty(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionUserRead)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-auth-stats").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("total_users").Number().IsEqual(1)     // 1 from GivenAuthToken
	obj.Value("total_roles").Number().IsEqual(0)     // no roles created
	obj.Value("active_sessions").Number().IsEqual(1) // 1 from GivenAuthToken
}

func TestGetAuthStats_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/auth/get-auth-stats")
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
