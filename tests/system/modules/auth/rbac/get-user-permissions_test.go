//go:build system

package rbac_test

import (
	"net/http"
	"testing"

	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"
)

func TestGetUserPermissions_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	auth.GivenUserPermissions(t,
		map[string]any{"user_id": u.ID, "permission": "perm1"},
		map[string]any{"user_id": u.ID, "permission": "perm2"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-user-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("user_id", u.ID).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().IsEqual(2)
}

func TestGetUserPermissions_UserNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-user-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("user_id", "00000000-0000-0000-0000-000000000000").
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestGetUserPermissions_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/auth/get-user-permissions")
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
