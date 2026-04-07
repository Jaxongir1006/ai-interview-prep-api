//go:build system

package taskmill_test

import (
	"net/http"
	"testing"
	"time"

	portalauth "go-enterprise-blueprint/internal/portal/auth"
	stateaudit "go-enterprise-blueprint/tests/state/audit"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/state/platform"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestCleanupResults_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillManage)
	now := time.Now()
	oldTime := now.Add(-48 * time.Hour)

	platform.GivenTaskResults(t,
		map[string]any{
			"queue_name":   "cleanup-test",
			"completed_at": oldTime,
		},
		map[string]any{
			"queue_name":   "cleanup-test",
			"completed_at": oldTime.Add(-1 * time.Hour),
		},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/cleanup-results").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"completed_before": now.Add(-24 * time.Hour).Format(time.RFC3339),
			"queue_name":       "cleanup-test",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("deleted_count").Number().Gt(0)

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}
