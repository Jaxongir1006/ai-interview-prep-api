//go:build system

package taskmill_test

import (
	"net/http"
	"testing"

	portalauth "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	stateaudit "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/platform"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestTriggerSchedule_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillManage)
	platform.GivenSchedules(t,
		map[string]any{
			"operation_id": "trigger-test-op",
			"queue_name":   "trigger-queue",
		},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/trigger-schedule").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"operation_id": "trigger-test-op",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestTriggerSchedule_NotFound(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillManage)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/trigger-schedule").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"operation_id": "non-existent-schedule",
		}).
		Expect()

	// THEN
	resp.Status(http.StatusNotFound)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("SCHEDULE_NOT_FOUND")
}
