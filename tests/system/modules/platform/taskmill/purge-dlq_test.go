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

func TestPurgeDLQ_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillManage)
	platform.GivenDLQTasks(t,
		map[string]any{"queue_name": "purge-dlq-test"},
		map[string]any{"queue_name": "purge-dlq-test"},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/purge-dlq").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{"queue_name": "purge-dlq-test"}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.Equal(t, 0, platform.GetDLQCount(t, "purge-dlq-test"),
		"DLQ should be empty after purge")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}
