//go:build system

package user_test

import (
	"net/http"
	"testing"
	"time"

	emailverificationdomain "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/emailverification"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	statecandidate "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestVerifyEmail_Success(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": false,
	})[0]

	rawToken := "raw-verification-token"
	tokenHash := emailverification.HashToken(rawToken)
	auth.GivenEmailVerificationTokens(t, map[string]any{
		"user_id":    u.ID,
		"email":      "candidate@example.com",
		"token_hash": tokenHash,
	})
	statecandidate.GivenProfiles(t, map[string]any{
		"user_id": u.ID,
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/verify-email").
		WithJSON(map[string]string{
			"token": rawToken,
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("user_id").String().IsEqual(u.ID)
	resp.JSON().Object().Value("email").String().IsEqual("candidate@example.com")
	resp.JSON().Object().Value("is_verified").Boolean().IsTrue()
	resp.JSON().Object().Value("access_token").String().NotEmpty()
	resp.JSON().Object().Value("access_token_expires_at").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token").String().NotEmpty()
	resp.JSON().Object().Value("refresh_token_expires_at").String().NotEmpty()
	resp.JSON().Object().Value("onboarding_required").Boolean().IsTrue()

	updatedUser := auth.GetUserByID(t, u.ID)
	assert.True(t, updatedUser.IsVerified)
	assert.Equal(t, 1, auth.SessionCount(t, u.ID))

	usedToken := auth.GetEmailVerificationTokenByHash(t, tokenHash)
	assert.NotNil(t, usedToken.UsedAt)
}

func TestVerifyEmail_OnboardingNotRequired(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": false,
	})[0]
	statecandidate.GivenProfiles(t, map[string]any{
		"user_id":                 u.ID,
		"target_role":             "python",
		"experience_level":        "junior",
		"onboarding_completed":    true,
		"onboarding_completed_at": time.Now().Add(-time.Hour),
	})

	rawToken := "raw-verification-token"
	auth.GivenEmailVerificationTokens(t, map[string]any{
		"user_id":    u.ID,
		"email":      "candidate@example.com",
		"token_hash": emailverification.HashToken(rawToken),
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/verify-email").
		WithJSON(map[string]string{
			"token": rawToken,
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("onboarding_required").Boolean().IsFalse()
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).POST("/api/v1/auth/verify-email").
		WithJSON(map[string]string{
			"token": "missing-token",
		}).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().Value("code").String().
		IsEqual(emailverificationdomain.CodeEmailVerificationTokenInvalid)
}

func TestVerifyEmail_ExpiredToken(t *testing.T) {
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{
		"email":       "candidate@example.com",
		"password":    auth.TestPassword1,
		"is_verified": false,
	})[0]

	rawToken := "expired-verification-token"
	auth.GivenEmailVerificationTokens(t, map[string]any{
		"user_id":    u.ID,
		"email":      "candidate@example.com",
		"token_hash": emailverification.HashToken(rawToken),
		"expires_at": time.Now().Add(-time.Hour),
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/verify-email").
		WithJSON(map[string]string{
			"token": rawToken,
		}).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().Value("code").String().
		IsEqual(emailverificationdomain.CodeEmailVerificationTokenExpired)
}
