//go:build system

package taskmill_test

import (
	"net/http"
	"testing"

	portalauth "go-enterprise-blueprint/internal/portal/auth"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/state/platform"
	"go-enterprise-blueprint/tests/system/trigger"
)

func TestListDLQTasks_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillView)
	platform.GivenDLQTasks(t,
		map[string]any{"queue_name": "dlq-queue", "operation_id": "test-op-1"},
		map[string]any{"queue_name": "dlq-queue", "operation_id": "test-op-2"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/list-dlq-tasks").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("queue_name", "dlq-queue").
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().Ge(2)
}
