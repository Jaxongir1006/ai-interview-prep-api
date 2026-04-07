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

func TestChangeMyPassword_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"password": auth.TestPassword1,
	})[0]
	sessions := auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})
	token := sessions[0].AccessToken

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/change-my-password").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{
			"current_password": auth.TestPassword1,
			"new_password":     auth.TestPassword2,
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)

	// Verify old password no longer works
	loginResp := trigger.UserAction(t).POST("/api/v1/auth/admin-login").
		WithJSON(map[string]string{
			"username": *u.Username,
			"password": auth.TestPassword1,
		}).
		Expect()
	loginResp.Status(http.StatusBadRequest)
	loginResp.JSON().Object().Value("error").Object().Value("code").String().
		IsEqual("INCORRECT_CREDENTIALS")

	// Verify new password works
	loginResp = trigger.UserAction(t).POST("/api/v1/auth/admin-login").
		WithJSON(map[string]string{
			"username": *u.Username,
			"password": auth.TestPassword2,
		}).
		Expect()
	loginResp.Status(http.StatusOK)

	// Verify audit log
	logs := stateaudit.GetActionLogs(t)
	// 2 logs: change-my-password + admin-login (the successful re-login)
	changeLog := logs[len(logs)-1]
	assert.Equal(t, "change-my-password", changeLog.OperationID)
	assert.NotNil(t, changeLog.UserID)
	assert.Equal(t, u.ID, *changeLog.UserID)
}

func TestChangeMyPassword_IncorrectCurrentPassword(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"password": auth.TestPassword1,
	})[0]
	sessions := auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})
	token := sessions[0].AccessToken

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/change-my-password").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{
			"current_password": auth.TestPassword2,
			"new_password":     auth.TestPassword3,
		}).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().Value("code").String().
		IsEqual("INCORRECT_CREDENTIALS")
	resp.JSON().Object().Value("error").Object().Value("message").String().
		NotContains("[untranslated]")
}

func TestChangeMyPassword_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "missing current_password",
			payload: map[string]string{"new_password": "newpassword123"},
		},
		{
			name:    "missing new_password",
			payload: map[string]string{"current_password": "oldpassword123"},
		},
		{
			name:    "new_password too short",
			payload: map[string]string{"current_password": "oldpassword123", "new_password": "short"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := auth.GivenAuthToken(t)

			// WHEN
			resp := trigger.UserAction(t).POST("/api/v1/auth/change-my-password").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(tc.payload).
				Expect()

			// THEN
			resp.Status(http.StatusBadRequest)
			resp.JSON().Object().Value("error").Object().Value("code").String().
				IsEqual("VALIDATION_FAILED")
			resp.JSON().Object().Value("error").Object().Value("message").String().
				NotContains("[untranslated]")
		})
	}
}
