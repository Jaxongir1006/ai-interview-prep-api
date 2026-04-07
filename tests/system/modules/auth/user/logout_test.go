//go:build system

package user_test

import (
	"net/http"
	"testing"

	stateaudit "github.com/jaxongir1006/hire-ready-api/tests/state/audit"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestLogout_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	sessions := auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})
	s := sessions[0]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/logout").
		WithHeader("Authorization", "Bearer "+s.AccessToken).
		WithJSON(map[string]any{}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.False(t, auth.SessionExists(t, s.ID),
		"session should be deleted after logout")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestLogout_Failures(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing authorization header",
			token:      "",
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name:       "invalid access token",
			token:      "nonexistent-token",
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)

			// WHEN
			req := trigger.UserAction(t).POST("/api/v1/auth/logout").
				WithJSON(map[string]any{})
			if tc.token != "" {
				req = req.WithHeader("Authorization", "Bearer "+tc.token)
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
