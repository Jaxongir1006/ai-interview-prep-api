//go:build system

package taskmill_test

import (
	"net/http"
	"testing"

	portalauth "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/platform"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"
)

func TestListTaskResults_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillView)
	platform.GivenTaskResults(t,
		map[string]any{"queue_name": "results-queue", "operation_id": "result-op-1"},
		map[string]any{"queue_name": "results-queue", "operation_id": "result-op-2"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/platform/list-task-results").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("queue_name", "results-queue").
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().Ge(2)
}
