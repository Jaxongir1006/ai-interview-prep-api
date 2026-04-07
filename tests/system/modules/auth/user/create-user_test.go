//go:build system

package user_test

import (
	"net/http"
	"strings"
	"testing"

	stateaudit "go-enterprise-blueprint/tests/state/audit"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/rise-and-shine/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/create-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{
			"username": "newuser",
			"password": "securepassword123",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("id").String().NotEmpty()
	resp.JSON().Object().Value("username").String().IsEqual("newuser")
	resp.JSON().Object().Value("is_active").Boolean().IsTrue()

	// Verify in DB
	u := auth.GetUserByUsername(t, "newuser")
	assert.True(t, u.IsActive)
	assert.True(t, strings.HasPrefix(*u.PasswordHash, "$2"), "password should be bcrypt hashed")
	assert.True(t, hasher.Compare("securepassword123", *u.PasswordHash))

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestCreateUser_UsernameConflict(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	auth.GivenUsers(t, map[string]any{"username": "existing"})

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/create-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]string{
			"username": "existing",
			"password": "securepassword123",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusConflict)
}

func TestCreateUser_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "missing username",
			payload: map[string]string{"password": "password123"},
		},
		{
			name:    "missing password",
			payload: map[string]string{"username": "testuser"},
		},
		{
			name:    "password too short",
			payload: map[string]string{"username": "testuser", "password": "short"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := auth.GivenSuperadminToken(t)

			// WHEN
			resp := trigger.UserAction(t).POST("/api/v1/auth/create-user").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(tc.payload).
				Expect()

			// THEN
			resp.Status(http.StatusBadRequest)
		})
	}
}

func TestCreateUser_AuthFailures(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing authorization header",
			setup:      func(_ *testing.T) string { return "" },
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name: "insufficient permissions",
			setup: func(t *testing.T) string {
				return auth.GivenAuthToken(t)
			},
			wantStatus: http.StatusForbidden,
			wantErr:    "FORBIDDEN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := tc.setup(t)

			// WHEN
			req := trigger.UserAction(t).POST("/api/v1/auth/create-user").
				WithJSON(map[string]any{})
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
