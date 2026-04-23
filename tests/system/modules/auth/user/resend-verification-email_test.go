//go:build system

package user_test

import (
	"net/http"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestResendVerificationEmail_Success(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": false,
	})[0]

	resp := trigger.UserAction(t).POST("/api/v1/auth/resend-verification-email").
		WithJSON(map[string]string{
			"email": "candidate@example.com",
		}).
		Expect()

	resp.Status(http.StatusOK)

	assert.True(t, auth.EmailVerificationTokenPointerExists(t, u.ID, "candidate@example.com"))
}

func TestResendVerificationEmail_MissingUserDoesNotLeak(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).POST("/api/v1/auth/resend-verification-email").
		WithJSON(map[string]string{
			"email": "missing@example.com",
		}).
		Expect()

	resp.Status(http.StatusOK)
}

func TestResendVerificationEmail_VerifiedUserDoesNotSend(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": true,
	})[0]

	resp := trigger.UserAction(t).POST("/api/v1/auth/resend-verification-email").
		WithJSON(map[string]string{
			"email": "candidate@example.com",
		}).
		Expect()

	resp.Status(http.StatusOK)
	assert.False(t, auth.EmailVerificationTokenPointerExists(t, u.ID, "candidate@example.com"))
}
