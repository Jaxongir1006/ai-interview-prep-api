//go:build system

package user_test

import (
	"net/http"
	"testing"

	passwordresetdomain "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/passwordresettoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/passwordreset"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/rise-and-shine/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestRequestPasswordReset_Success(t *testing.T) {
	database.Empty(t)
	auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": true,
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/request-password-reset").
		WithJSON(map[string]string{
			"email": "candidate@example.com",
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("message").String().NotEmpty()
}

func TestRequestPasswordReset_MissingUserDoesNotLeak(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).POST("/api/v1/auth/request-password-reset").
		WithJSON(map[string]string{
			"email": "missing@example.com",
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("message").String().NotEmpty()
}

func TestConfirmPasswordReset_Success(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": false,
	})[0]
	auth.GivenSessions(t, map[string]any{
		"user_id": u.ID,
	})

	rawToken := "raw-password-reset-token"
	auth.GivenPasswordResetTokens(t, map[string]any{
		"user_id":    u.ID,
		"email":      "candidate@example.com",
		"token_hash": passwordreset.HashToken(rawToken),
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/confirm-password-reset").
		WithJSON(map[string]string{
			"token":    rawToken,
			"password": "NewSecurePassword_1",
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("message").String().NotEmpty()

	updatedUser := auth.GetUserByID(t, u.ID)
	assert.True(t, updatedUser.IsVerified)
	assert.True(t, hasher.Compare("NewSecurePassword_1", *updatedUser.PasswordHash))
	assert.Equal(t, 0, auth.SessionCount(t, u.ID))
}

func TestConfirmPasswordReset_InvalidToken(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).POST("/api/v1/auth/confirm-password-reset").
		WithJSON(map[string]string{
			"token":    "missing-token",
			"password": "NewSecurePassword_1",
		}).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().Value("code").String().
		IsEqual(passwordresetdomain.CodePasswordResetTokenInvalid)
}

func TestConfirmPasswordReset_UsedToken(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": true,
	})[0]

	rawToken := "raw-password-reset-token"
	auth.GivenPasswordResetTokens(t, map[string]any{
		"user_id":    u.ID,
		"email":      "candidate@example.com",
		"token_hash": passwordreset.HashToken(rawToken),
	})

	trigger.UserAction(t).POST("/api/v1/auth/confirm-password-reset").
		WithJSON(map[string]string{
			"token":    rawToken,
			"password": "NewSecurePassword_1",
		}).
		Expect().
		Status(http.StatusOK)

	resp := trigger.UserAction(t).POST("/api/v1/auth/confirm-password-reset").
		WithJSON(map[string]string{
			"token":    rawToken,
			"password": "AnotherSecurePassword_1",
		}).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().Value("code").String().
		IsEqual(passwordresetdomain.CodePasswordResetTokenInvalid)
}
