//go:build system

package alerterror_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	stateaudit "github.com/jaxongir1006/hire-ready-api/tests/state/audit"
	stateauth "github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/state/platform"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestCleanupErrors_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := stateauth.GivenAuthToken(t, auth.PermissionAlertManage)
	now := time.Now()

	platform.GivenErrors(t,
		map[string]any{"created_at": now.Add(-48 * time.Hour)},
		map[string]any{"created_at": now.Add(-72 * time.Hour)},
		map[string]any{"created_at": now.Add(-1 * time.Hour)}, // recent, should NOT be deleted
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/cleanup-errors").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"created_before": now.Add(-24 * time.Hour).Format(time.RFC3339),
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("deleted_count").Number().IsEqual(2)

	// Verify remaining count
	remaining := platform.GetErrorCount(t)
	if remaining != 1 {
		t.Errorf("expected 1 remaining error, got %d", remaining)
	}

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestCleanupErrors_AuthFailures(t *testing.T) {
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
				return stateauth.GivenAuthToken(t)
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
			now := time.Now()

			// WHEN
			req := trigger.UserAction(t).POST("/api/v1/platform/cleanup-errors").
				WithJSON(map[string]any{
					"created_before": now.Add(-24 * time.Hour).Format(time.RFC3339),
				})
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
