//go:build system

package session_test

import (
	"net/http"
	"testing"

	stateaudit "github.com/jaxongir1006/hire-ready-api/tests/state/audit"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestDeleteMySession_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	sessions := auth.GivenUserSessions(t, u.ID, 2)
	currentSession := sessions[0]
	targetSession := sessions[1]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-my-session").
		WithHeader("Authorization", "Bearer "+currentSession.AccessToken).
		WithJSON(map[string]any{
			"session_id": targetSession.ID,
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.False(t, auth.SessionExists(t, targetSession.ID),
		"target session should be deleted")
	assert.True(t, auth.SessionExists(t, currentSession.ID),
		"current session should still exist")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestDeleteMySession_Failures(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) (string, map[string]any)
		wantStatus int
		wantErr    string
	}{
		{
			name: "missing authorization header",
			setup: func(_ *testing.T) (string, map[string]any) {
				return "", map[string]any{"session_id": int64(999)}
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name: "session not found",
			setup: func(t *testing.T) (string, map[string]any) {
				u := auth.GivenUsers(t, map[string]any{})[0]
				sessions := auth.GivenSessions(t, map[string]any{
					"user_id": u.ID,
				})
				return sessions[0].AccessToken, map[string]any{
					"session_id": int64(999999),
				}
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "SESSION_NOT_FOUND",
		},
		{
			name: "session belongs to another user",
			setup: func(t *testing.T) (string, map[string]any) {
				users := auth.GivenUsers(t,
					map[string]any{"username": "user1"},
					map[string]any{"username": "user2"},
				)
				user1Sessions := auth.GivenSessions(t, map[string]any{
					"user_id": users[0].ID,
				})
				user2Sessions := auth.GivenSessions(t, map[string]any{
					"user_id": users[1].ID,
				})
				return user1Sessions[0].AccessToken, map[string]any{
					"session_id": user2Sessions[0].ID,
				}
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "SESSION_NOT_FOUND",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token, payload := tc.setup(t)

			// WHEN
			req := trigger.UserAction(t).POST("/api/v1/auth/delete-my-session").
				WithJSON(payload)
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

func TestDeleteMySession_ValidationErrors(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	sessions := auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/delete-my-session").
		WithHeader("Authorization", "Bearer "+sessions[0].AccessToken).
		WithJSON(map[string]any{}).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("VALIDATION_FAILED")
	resp.JSON().Object().Value("error").Object().
		Value("message").String().NotContains("[untranslated]")
}
