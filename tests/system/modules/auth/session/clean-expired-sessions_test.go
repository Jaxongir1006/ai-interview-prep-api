//go:build system

package session_test

import (
	"testing"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/usecase/session/cleanexpiredsessions"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestCleanExpiredSessions(t *testing.T) {
	// GIVEN - sessions with various expiry states
	database.Empty(t)
	u := auth.GivenUsers(t, map[string]any{})[0]

	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(24 * time.Hour)

	sessions := auth.GivenSessions(t,
		// Should be DELETED - both tokens expired
		map[string]any{
			"user_id":                  u.ID,
			"access_token_expires_at":  past,
			"refresh_token_expires_at": past,
		},
		// Should be DELETED - refresh expired (access token status irrelevant)
		map[string]any{
			"user_id":                  u.ID,
			"access_token_expires_at":  future,
			"refresh_token_expires_at": past,
		},
		// Should be PRESERVED - refresh token valid (even with expired access)
		map[string]any{
			"user_id":                  u.ID,
			"access_token_expires_at":  past,
			"refresh_token_expires_at": future,
		},
		// Should be PRESERVED - both tokens valid
		map[string]any{
			"user_id":                  u.ID,
			"access_token_expires_at":  future,
			"refresh_token_expires_at": future,
		},
	)

	// WHEN
	trigger.AsyncTask(t, "auth", "clean-expired-sessions", &cleanexpiredsessions.Payload{})

	// THEN
	assert.False(t, auth.SessionExists(t, sessions[0].ID), "session with both tokens expired should be deleted")
	assert.False(t, auth.SessionExists(t, sessions[1].ID), "session with expired refresh should be deleted")
	assert.True(t, auth.SessionExists(t, sessions[2].ID), "session with valid refresh should be preserved")
	assert.True(t, auth.SessionExists(t, sessions[3].ID), "session with both tokens valid should be preserved")
}
