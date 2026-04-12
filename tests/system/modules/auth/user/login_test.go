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

func TestLogin_Success(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": true,
	})[0]

	sessionCountBefore := auth.SessionCount(t, u.ID)

	resp := trigger.UserAction(t).POST("/api/v1/auth/login").
		WithJSON(map[string]string{
			"email":    "candidate@example.com",
			"password": auth.TestPassword1,
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("access_token").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token").String().NotEmpty()

	assert.Equal(t, sessionCountBefore+1, auth.SessionCount(t, u.ID))

	updatedUser := auth.GetUserByID(t, u.ID)
	assert.NotNil(t, updatedUser.LastLoginAt)
	assert.NotNil(t, updatedUser.LastActiveAt)
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestLogin_IncorrectCredentials(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) map[string]string
		wantCode string
	}{
		{
			name: "user not found",
			setup: func(_ *testing.T) map[string]string {
				return map[string]string{
					"email":    "missing@example.com",
					"password": auth.TestPassword1,
				}
			},
			wantCode: "INCORRECT_CREDENTIALS",
		},
		{
			name: "incorrect password",
			setup: func(t *testing.T) map[string]string {
				auth.GivenUsers(t, map[string]any{
					"email":       "candidate@example.com",
					"password":    auth.TestPassword1,
					"is_verified": true,
				})
				return map[string]string{
					"email":    "candidate@example.com",
					"password": auth.TestPassword2,
				}
			},
			wantCode: "INCORRECT_CREDENTIALS",
		},
		{
			name: "user inactive",
			setup: func(t *testing.T) map[string]string {
				auth.GivenUsers(t, map[string]any{
					"email":     "candidate@example.com",
					"password":  auth.TestPassword1,
					"is_active": false,
				})
				return map[string]string{
					"email":    "candidate@example.com",
					"password": auth.TestPassword1,
				}
			},
			wantCode: "INCORRECT_CREDENTIALS",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			database.Empty(t)

			resp := trigger.UserAction(t).POST("/api/v1/auth/login").
				WithJSON(tc.setup(t)).
				Expect()

			resp.Status(http.StatusBadRequest)
			resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual(tc.wantCode)
		})
	}
}

func TestLogin_EmailNotVerified(t *testing.T) {
	database.Empty(t)
	auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": false,
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/login").
		WithJSON(map[string]string{
			"email":    "candidate@example.com",
			"password": auth.TestPassword1,
		}).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual("EMAIL_NOT_VERIFIED")
}
