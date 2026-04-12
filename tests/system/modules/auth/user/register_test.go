//go:build system

package user_test

import (
	"net/http"
	"strings"
	"testing"

	stateaudit "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	statecandidate "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/rise-and-shine/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestRegister_CORSPreflight(t *testing.T) {
	resp := trigger.UserAction(t).OPTIONS("/api/v1/auth/register").
		WithHeader("Origin", "http://localhost:3000").
		WithHeader("Access-Control-Request-Method", http.MethodPost).
		WithHeader("Access-Control-Request-Headers", "content-type,authorization").
		Expect()

	resp.Status(http.StatusNoContent)
	resp.Header("Access-Control-Allow-Origin").IsEqual("http://localhost:3000")
	resp.Header("Access-Control-Allow-Credentials").IsEqual("true")
	resp.Header("Access-Control-Allow-Methods").Contains(http.MethodPost)
	resp.Header("Access-Control-Allow-Methods").Contains(http.MethodOptions)
	resp.Header("Access-Control-Allow-Headers").Contains("Content-Type")
	resp.Header("Access-Control-Allow-Headers").Contains("Authorization")
	resp.Header("Access-Control-Max-Age").IsEqual("86400")
}

func TestRegister_Success(t *testing.T) {
	database.Empty(t)

	resp := trigger.UserAction(t).POST("/api/v1/auth/register").
		WithJSON(map[string]string{
			"email":     "candidate@example.com",
			"full_name": "John Candidate",
			"password":  "SecurePassword_1",
		}).
		Expect()

	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("user").Object().Value("id").String().NotEmpty()
	resp.JSON().Object().Value("user").Object().Value("email").String().IsEqual("candidate@example.com")
	resp.JSON().Object().Value("user").Object().Value("is_verified").Boolean().IsFalse()
	resp.JSON().Object().Value("profile").Object().Value("full_name").String().IsEqual("John Candidate")
	resp.JSON().Object().Value("profile").Object().Value("preferred_topics").Array().Length().IsEqual(0)
	resp.JSON().Object().Value("verification_required").Boolean().IsTrue()

	u := auth.GetUserByEmail(t, "candidate@example.com")
	assert.True(t, u.IsActive)
	assert.NotNil(t, u.PasswordHash)
	assert.True(t, strings.HasPrefix(*u.PasswordHash, "$2"))
	assert.True(t, hasher.Compare("SecurePassword_1", *u.PasswordHash))
	assert.False(t, u.IsVerified)

	verificationTokens := auth.GetEmailVerificationTokensByUserID(t, u.ID)
	assert.Len(t, verificationTokens, 1)
	assert.Equal(t, "candidate@example.com", verificationTokens[0].Email)
	assert.NotEmpty(t, verificationTokens[0].TokenHash)
	assert.Nil(t, verificationTokens[0].UsedAt)

	profile := statecandidate.GetProfileByUserID(t, u.ID)
	assert.NotNil(t, profile.FullName)
	assert.Equal(t, "John Candidate", *profile.FullName)
	assert.Nil(t, profile.TargetRole)
	assert.Nil(t, profile.ExperienceLevel)
	assert.Equal(t, 3, profile.InterviewGoalPerWeek)

	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestRegister_EmailConflict(t *testing.T) {
	database.Empty(t)
	auth.GivenUsers(t, map[string]any{
		"email":    "candidate@example.com",
		"password": auth.TestPassword1,
	})

	resp := trigger.UserAction(t).POST("/api/v1/auth/register").
		WithJSON(map[string]string{
			"email":     "candidate@example.com",
			"full_name": "John Candidate",
			"password":  "SecurePassword_1",
		}).
		Expect()

	resp.Status(http.StatusConflict)
	resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual("EMAIL_CONFLICT")
}

func TestRegister_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]string
	}{
		{
			name: "missing email",
			payload: map[string]string{
				"full_name": "John Candidate",
				"password":  "SecurePassword_1",
			},
		},
		{
			name: "missing full_name",
			payload: map[string]string{
				"email":    "candidate@example.com",
				"password": "SecurePassword_1",
			},
		},
		{
			name: "missing password",
			payload: map[string]string{
				"email":     "candidate@example.com",
				"full_name": "John Candidate",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			database.Empty(t)

			resp := trigger.UserAction(t).POST("/api/v1/auth/register").
				WithJSON(tc.payload).
				Expect()

			resp.Status(http.StatusBadRequest)
			resp.JSON().Object().Value("error").Object().Value("code").String().IsEqual("VALIDATION_FAILED")
		})
	}
}
