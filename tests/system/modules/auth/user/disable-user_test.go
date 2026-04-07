//go:build system

package user_test

import (
	"net/http"
	"testing"

	stateaudit "go-enterprise-blueprint/tests/state/audit"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestDisableUser_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	auth.GivenSessions(t,
		map[string]any{"user_id": u.ID},
		map[string]any{"user_id": u.ID},
	)
	assert.Equal(t, 2, auth.SessionCount(t, u.ID))

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/disable-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{"id": u.ID}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)

	updated := auth.GetUserByID(t, u.ID)
	assert.False(t, updated.IsActive, "user should be disabled")
	assert.Equal(t, 0, auth.SessionCount(t, u.ID), "all sessions should be deleted")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestDisableUser_AlreadyDisabled(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{"is_active": false})[0]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/disable-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{"id": u.ID}).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
}

func TestDisableUser_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/disable-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{"id": "00000000-0000-0000-0000-000000000000"}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestDisableUser_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/disable-user").
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
