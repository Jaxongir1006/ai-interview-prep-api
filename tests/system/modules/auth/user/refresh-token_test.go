//go:build system

package user_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	sessions := auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})
	originalSession := sessions[0]

	sessionCountBefore := auth.SessionCount(t, u.ID)
	assert.Nil(t, u.LastActiveAt, "last_active_at should be nil initially")

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/refresh-token").
		WithJSON(map[string]string{
			"refresh_token": originalSession.RefreshToken,
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("access_token").String().NotEmpty()
	resp.JSON().Object().Value("access_token_expires_at").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token_expires_at").String().NotEmpty()

	// Session count should remain the same (updated, not created)
	sessionCountAfter := auth.SessionCount(t, u.ID)
	assert.Equal(t, sessionCountBefore, sessionCountAfter,
		"session count should not change")

	// User's last_active_at should be updated
	updatedUser := auth.GetUserByID(t, u.ID)
	assert.NotNil(t, updatedUser.LastActiveAt,
		"last_active_at should be updated")
}

func TestRefreshToken_Failures(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr string
	}{
		{
			name: "invalid refresh token",
			setup: func(t *testing.T) string {
				return "nonexistent-token"
			},
			wantErr: "SESSION_NOT_FOUND",
		},
		{
			name: "expired refresh token",
			setup: func(t *testing.T) string {
				u := auth.GivenUsers(t, map[string]any{})[0]
				sessions := auth.GivenSessions(t, map[string]any{
					"user_id":                  u.ID,
					"refresh_token_expires_at": time.Now().Add(-1 * time.Hour),
				})
				return sessions[0].RefreshToken
			},
			wantErr: "SESSION_NOT_FOUND",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			refreshToken := tc.setup(t)

			// WHEN
			resp := trigger.UserAction(t).POST("/api/v1/auth/refresh-token").
				WithJSON(map[string]string{
					"refresh_token": refreshToken,
				}).
				Expect()

			// THEN
			resp.Status(http.StatusBadRequest)
			resp.JSON().Object().Value("error").Object().
				Value("code").String().IsEqual(tc.wantErr)
			resp.JSON().Object().Value("error").Object().
				Value("message").String().NotContains("[untranslated]")
		})
	}
}

func TestRefreshToken_ValidationErrors(t *testing.T) {
	// GIVEN
	database.Empty(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/refresh-token").
		WithJSON(map[string]string{}).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("VALIDATION_FAILED")
	resp.JSON().Object().Value("error").Object().
		Value("message").String().NotContains("[untranslated]")
}
