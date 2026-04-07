//go:build system

package taskmill_test

import (
	"net/http"
	"testing"

	portalauth "go-enterprise-blueprint/internal/portal/auth"
	stateaudit "go-enterprise-blueprint/tests/state/audit"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/state/platform"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestPurgeQueue_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillManage)
	platform.GivenQueuedTasks(t,
		map[string]any{"queue_name": "purge-test"},
		map[string]any{"queue_name": "purge-test"},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/purge-queue").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"queue_name": "purge-test"}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.Equal(t, 0, platform.GetTaskQueueCount(t, "purge-test"),
		"queue should be empty after purge")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}
