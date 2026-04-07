//go:build system

package rbac_test

import (
	"net/http"
	"testing"

	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"
)

func TestGetRolePermissions_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	role := auth.GivenRoles(t, map[string]any{"name": "admin"})[0]
	auth.GivenRolePermissions(t,
		map[string]any{"role_id": role.ID, "permission": "perm1"},
		map[string]any{"role_id": role.ID, "permission": "perm2"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-role-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("role_id", role.ID).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().IsEqual(2)
}

func TestGetRolePermissions_RoleNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-role-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("role_id", 99999).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestGetRolePermissions_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).GET("/api/v1/auth/get-role-permissions")
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
