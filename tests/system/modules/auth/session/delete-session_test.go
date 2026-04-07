//go:build system

package session_test

import (
	"net/http"
	"testing"

	stateaudit "go-enterprise-blueprint/tests/state/audit"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestDeleteSession_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	sessions := auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})
	s := sessions[0]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-session").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"session_id": s.ID}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.False(t, auth.SessionExists(t, s.ID), "session should be deleted")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestDeleteSession_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-session").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"session_id": 99999}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestDeleteSession_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/delete-session").
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
