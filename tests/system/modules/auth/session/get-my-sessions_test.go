//go:build system

package session_test

import (
	"net/http"
	"testing"

	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"
)

func TestGetMySessions_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{})[0]
	auth.GivenUserSessions(t, u.ID, 3)

	sessions := auth.GetSessionsByUserID(t, u.ID)
	accessToken := sessions[0].AccessToken

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-my-sessions").
		WithHeader("Authorization", "Bearer "+accessToken).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().IsEqual(3)

	// Verify session fields are present (tokens excluded via json:"-")
	first := resp.JSON().Object().Value("content").Array().Value(0).Object()
	first.Value("id").Number().Gt(0)
	first.Value("ip_address").String().NotEmpty()
	first.Value("user_agent").String().NotEmpty()
	first.Value("last_used_at").String().NotEmpty()
}

func TestGetMySessions_Failures(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing authorization header",
			token:      "",
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name:       "invalid access token",
			token:      "nonexistent-token",
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)

			// WHEN
			req := trigger.UserAction(t).GET("/api/v1/auth/get-my-sessions")
			if tc.token != "" {
				req = req.WithHeader("Authorization", "Bearer "+tc.token)
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
