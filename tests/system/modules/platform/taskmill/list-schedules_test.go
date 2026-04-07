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

func TestListSchedules_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillView)
	platform.GivenSchedules(t,
		map[string]any{"operation_id": "schedule-op-1", "queue_name": "sched-queue"},
		map[string]any{"operation_id": "schedule-op-2", "queue_name": "sched-queue"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/list-schedules").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().Ge(2)
}
