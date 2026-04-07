//go:build system

package user_test

import (
	"net/http"
	"testing"

	stateaudit "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/rise-and-shine/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	u := auth.GivenUsers(t, map[string]any{
		"username": "oldname",
		"password": auth.TestPassword1,
	})[0]

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/update-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"id":       u.ID,
			"username": "newname",
			"password": "newsecurepassword",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("username").String().IsEqual("newname")

	// Verify in DB
	updated := auth.GetUserByID(t, u.ID)
	assert.Equal(t, "newname", *updated.Username)
	assert.True(t, hasher.Compare("newsecurepassword", *updated.PasswordHash))

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestUpdateUser_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/auth/update-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"id":       "00000000-0000-0000-0000-000000000000",
			"username": "newname",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
}

func TestUpdateUser_UsernameConflict(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	users := auth.GivenUsers(t,
		map[string]any{"username": "user1"},
		map[string]any{"username": "user2"},
	)

	// WHEN - try to rename user2 to user1
	resp := trigger.UserAction(t).POST("/api/v1/auth/update-user").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"id":       users[1].ID,
			"username": "user1",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusConflict)
}

func TestUpdateUser_AuthFailures(t *testing.T) {
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
			req := trigger.UserAction(t).POST("/api/v1/auth/update-user").
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
