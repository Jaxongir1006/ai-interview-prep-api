//go:build system

package taskmill_test

import (
	"net/http"
	"testing"

	portalauth "github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	stateaudit "github.com/jaxongir1006/hire-ready-api/tests/state/audit"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/state/platform"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestRequeueFromDLQ_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalauth.PermissionTaskmillManage)
	ids := platform.GivenDLQTasks(t,
		map[string]any{"queue_name": "requeue-test", "operation_id": "requeue-op"},
	)

	// WHEN
	resp := trigger.UserAction(t).POST("/api/v1/platform/requeue-from-dlq").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]any{
			"task_id": ids[0],
		}).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	assert.Equal(t, 0, platform.GetDLQCount(t, "requeue-test"),
		"DLQ should be empty after requeue")
	assert.Equal(t, 1, platform.GetTaskQueueCount(t, "requeue-test"),
		"task should be back in the queue")

	// Verify audit log
	assert.Equal(t, 1, stateaudit.ActionLogCount(t))
}

func TestRequeueFromDLQ_AuthFailures(t *testing.T) {
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
				// User has taskmill:view but needs taskmill:manage
				return auth.GivenAuthToken(t, portalauth.PermissionTaskmillView)
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

			// WHEN
			req := trigger.UserAction(t).POST("/api/v1/platform/requeue-from-dlq").
				WithJSON(map[string]any{})
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
