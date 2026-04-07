//go:build system

package rbac_test

import (
	"net/http"
	"testing"

	stateaudit "github.com/jaxongir1006/hire-ready-api/tests/state/audit"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestSetUserPermissions_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	auth.GivenUserPermissions(t,
		map[string]any{"user_id": u.ID, "permission": "old:perm"},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-user-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"user_id":     u.ID,
			"permissions": []string{"new:perm1", "new:perm2"},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)

	perms := auth.GetUserPermissions(t, u.ID)
	assert.Len(t, perms, 2, "user should have exactly 2 direct permissions")

	permNames := make([]string, len(perms))
	for i, p := range perms {
		permNames[i] = p.Permission
	}
	assert.Contains(t, permNames, "new:perm1")
	assert.Contains(t, permNames, "new:perm2")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestSetUserPermissions_UserNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-user-permissions").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"user_id":     "00000000-0000-0000-0000-000000000000",
			"permissions": []string{"perm1"},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestSetUserPermissions_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/set-user-permissions").
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
