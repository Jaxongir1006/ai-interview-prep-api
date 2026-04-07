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

func TestGetQueueStats_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillView)
	platform.GivenQueuedTasks(t,
		map[string]any{"queue_name": "stats-queue"},
		map[string]any{"queue_name": "stats-queue"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/get-queue-stats").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("queue_name", "stats-queue").
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().NotEmpty()
	resp.JSON().Object().ContainsKey("total")
	resp.JSON().Object().ContainsKey("available")
}

func TestGetQueueStats_ValidationFailure(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillView)

	// WHEN (missing required queue_name parameter)
	resp := trigger.UserAction(t).GET("/api/v1/platform/get-queue-stats").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("VALIDATION_FAILED")
}
