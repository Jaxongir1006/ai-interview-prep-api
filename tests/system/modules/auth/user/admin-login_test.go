//go:build system

package user_test

import (
	"net/http"
	"testing"

	stateaudit "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestAdminLogin_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	users := auth.GivenUsers(t, map[string]any{
		"username": "testadmin",
		"password": auth.TestPassword1,
	})
	u := users[0]

	sessionCountBefore := auth.SessionCount(t, u.ID)
	assert.Nil(t, u.LastActiveAt, "last_active_at should be nil initially")

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/admin-login").
		WithJSON(map[string]string{
			"username": *u.Username,
			"password": auth.TestPassword1,
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("access_token").String().NotEmpty()
	resp.JSON().Object().Value("access_token_expires_at").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token_expires_at").String().NotEmpty()

	// Session created
	sessionCountAfter := auth.SessionCount(t, u.ID)
	assert.Equal(t, sessionCountBefore+1, sessionCountAfter, "one new session should be created")

	// LastActiveAt updated
	updatedUser := auth.GetUserByID(t, u.ID)
	assert.NotNil(t, updatedUser.LastActiveAt, "last_active_at should be updated after login")

	// Verify audit log with correct user_id
	logs := stateaudit.GetActionLogs(t)
	assert.Len(t, logs, 1)
	assert.NotNil(t, logs[0].UserID)
	assert.Equal(t, u.ID, *logs[0].UserID)
	assert.Equal(t, "admin-login", logs[0].OperationID)
}

func TestAdminLogin_IncorrectCredentials(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (username, password string)
		wantCode string
	}{
		{
			name: "user not found",
			setup: func(t *testing.T) (string, string) {
				return "nonexistent", "anypassword"
			},
			wantCode: "INCORRECT_CREDENTIALS",
		},
		{
			name: "incorrect password",
			setup: func(t *testing.T) (string, string) {
				users := auth.GivenUsers(t, map[string]any{
					"password": auth.TestPassword1,
				})
				return *users[0].Username, auth.TestPassword2
			},
			wantCode: "INCORRECT_CREDENTIALS",
		},
		{
			name: "user inactive",
			setup: func(t *testing.T) (string, string) {
				users := auth.GivenUsers(t, map[string]any{
					"is_active": false,
				})
				return *users[0].Username, auth.TestPassword1
			},
			wantCode: "INCORRECT_CREDENTIALS",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			username, password := tc.setup(t)

			// WHEN
			resp := trigger.UserAction(t).POST("/api/v1/auth/admin-login").
				WithJSON(map[string]string{
					"username": username,
					"password": password,
				}).
				Expect()

			// THEN
			resp.Status(http.StatusBadRequest)
			resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual(tc.wantCode)
			resp.JSON().Object().Value("error").Object().Value("message").String().NotContains("[untranslated]")
		})
	}
}

func TestAdminLogin_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "missing username",
			payload: map[string]string{"password": "somepassword"},
		},
		{
			name:    "missing password",
			payload: map[string]string{"username": "someuser"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)

			// WHEN
			resp := trigger.UserAction(t).POST("/api/v1/auth/admin-login").
				WithJSON(tc.payload).
				Expect()

			// THEN
			resp.Status(http.StatusBadRequest)
			resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual("VALIDATION_FAILED")
			resp.JSON().Object().Value("error").Object().Value("message").String().NotContains("[untranslated]")
		})
	}
}

func TestAdminLogin_MaxSessionsEnforcement(t *testing.T) {
	// GIVEN
	database.Empty(t)
	users := auth.GivenUsers(t, map[string]any{})
	u := users[0]

	maxSessions := 5 // should be consistent with test.config -> auth.max_active_sessions
	existingSessions := auth.GivenUserSessions(t, u.ID, maxSessions)
	oldestSession := existingSessions[0]

	assert.Equal(t, maxSessions, auth.SessionCount(t, u.ID))

	// WHEN
	trigger.UserAction(t).POST("/api/v1/auth/admin-login").
		WithJSON(map[string]string{
			"username": *u.Username,
			"password": auth.TestPassword1,
		}).
		Expect().
		Status(http.StatusOK)

	// THEN
	assert.Equal(t, maxSessions, auth.SessionCount(t, u.ID),
		"session count should remain at max after login")
	assert.False(t, auth.SessionExists(t, oldestSession.ID),
		"oldest session should be deleted")
}
