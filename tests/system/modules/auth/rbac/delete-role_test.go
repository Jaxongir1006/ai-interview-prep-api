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

func TestDeleteRole_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	role := auth.GivenRoles(t, map[string]any{"name": "to-delete"})[0]
	auth.GivenRolePermissions(t,
		map[string]any{"role_id": role.ID, "permission": "perm1"},
		map[string]any{"role_id": role.ID, "permission": "perm2"},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-role").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"id": role.ID}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.False(t, auth.RoleExists(t, role.ID), "role should be deleted")
	assert.Empty(t, auth.GetRolePermissions(t, role.ID), "role permissions should be deleted")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestDeleteRole_HasAssignedUsers(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	role := auth.GivenRoles(t, map[string]any{"name": "assigned"})[0]
	auth.GivenUserRoles(t, map[string]any{
		"user_id": u.ID,
		"role_id": role.ID,
	})

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-role").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"id": role.ID}).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
	assert.True(t, auth.RoleExists(t, role.ID), "role should not be deleted")
}

func TestDeleteRole_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-role").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"id": 99999}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestDeleteRole_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/delete-role").
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
