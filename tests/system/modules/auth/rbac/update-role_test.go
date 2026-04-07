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

func TestUpdateRole_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	role := auth.GivenRoles(t, map[string]any{"name": "old-name"})[0]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/update-role").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"id":   role.ID,
			"name": "new-name",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("name").String().IsEqual("new-name")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestUpdateRole_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/update-role").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"id": 99999, "name": "test"}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestUpdateRole_NameConflict(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	roles := auth.GivenRoles(t,
		map[string]any{"name": "role1"},
		map[string]any{"name": "role2"},
	)

	// WHEN - rename role2 to role1
	resp := trigger.UserAction(t).POST("/api/v1/auth/update-role").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"id":   roles[1].ID,
			"name": "role1",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusConflict)
}

func TestUpdateRole_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/update-role").
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
