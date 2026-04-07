//go:build system

package rbac_test

import (
	"net/http"
	"testing"

	stateaudit "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestSetRolePermissions_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	role := auth.GivenRoles(t, map[string]any{"name": "admin"})[0]
	auth.GivenRolePermissions(t,
		map[string]any{"role_id": role.ID, "permission": "old:perm"},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-role-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"role_id":     role.ID,
			"permissions": []string{"new:perm1", "new:perm2"},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)

	perms := auth.GetRolePermissions(t, role.ID)
	assert.Len(t, perms, 2, "should have exactly 2 permissions")

	permNames := make([]string, len(perms))
	for i, p := range perms {
		permNames[i] = p.Permission
	}
	assert.Contains(t, permNames, "new:perm1")
	assert.Contains(t, permNames, "new:perm2")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestSetRolePermissions_RoleNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-role-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"role_id":     99999,
			"permissions": []string{"perm1"},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestSetRolePermissions_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/set-role-permissions").
				WithJSON(map[string]any{})
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
