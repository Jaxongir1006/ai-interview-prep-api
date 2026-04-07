//go:build system

package taskmill_test

import (
	"net/http"
	"testing"

	portalauth "github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/state/platform"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"
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
