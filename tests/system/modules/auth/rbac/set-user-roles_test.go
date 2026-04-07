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

func TestSetUserRoles_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	roles := auth.GivenRoles(t,
		map[string]any{"name": "admin"},
		map[string]any{"name": "editor"},
	)
	// Pre-assign one role to verify replacement
	auth.GivenUserRoles(t, map[string]any{
		"user_id": u.ID,
		"role_id": roles[0].ID,
	})

	// WHEN - replace with both roles
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-user-roles").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"user_id":  u.ID,
			"role_ids": []int64{roles[0].ID, roles[1].ID},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)

	userRoles := auth.GetUserRoles(t, u.ID)
	assert.Len(t, userRoles, 2, "user should have exactly 2 role assignments")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestSetUserRoles_UserNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-user-roles").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"user_id":  "00000000-0000-0000-0000-000000000000",
			"role_ids": []int64{1},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestSetUserRoles_RoleNotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/set-user-roles").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"user_id":  u.ID,
			"role_ids": []int64{99999},
		}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestSetUserRoles_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/set-user-roles").
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
