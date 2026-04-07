//go:build system

package actionlog_test

import (
	"net/http"
	"testing"
	"time"

	portalaudit "github.com/jaxongir1006/hire-ready-api/internal/portal/audit"
	"github.com/jaxongir1006/hire-ready-api/tests/state/audit"
	"github.com/jaxongir1006/hire-ready-api/tests/state/auth"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
	"github.com/jaxongir1006/hire-ready-api/tests/system/trigger"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func TestGetActionLogs_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionActionLogRead)

	now := time.Now()
	audit.GivenActionLogs(t,
		map[string]any{
			"module":       "auth",
			"operation_id": "admin-login",
			"ip_address":   "10.0.0.1",
			"user_agent":   "Mozilla/5.0",
			"trace_id":     "trace-1",
			"created_at":   now.Add(-2 * time.Hour),
		},
		map[string]any{
			"module":       "auth",
			"operation_id": "create-user",
			"ip_address":   "10.0.0.2",
			"user_agent":   "curl/7.0",
			"trace_id":     "trace-2",
			"created_at":   now.Add(-1 * time.Hour),
		},
		map[string]any{
			"module":       "platform",
			"operation_id": "list-queues",
			"trace_id":     "trace-3",
			"created_at":   now,
		},
	)

	from := now.Add(-3 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
		WithQuery("from", from).
		WithQuery("to", to).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	arr := resp.JSON().Array()
	arr.Length().IsEqual(3)

	// Results should be in id DESC order (newest first)
	arr.Value(0).Object().Value("module").String().IsEqual("platform")
	arr.Value(1).Object().Value("operation_id").String().IsEqual("create-user")
	arr.Value(2).Object().Value("operation_id").String().IsEqual("admin-login")

	// Verify all fields are present
	first := arr.Value(0).Object()
	first.ContainsKey("id")
	first.ContainsKey("user_id")
	first.ContainsKey("module")
	first.ContainsKey("operation_id")
	first.ContainsKey("request_payload")
	first.ContainsKey("tags")
	first.ContainsKey("group_key")
	first.ContainsKey("ip_address")
	first.ContainsKey("user_agent")
	first.ContainsKey("trace_id")
	first.ContainsKey("created_at")
}

func TestGetActionLogs_FilterByModule(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionActionLogRead)

	now := time.Now()
	audit.GivenActionLogs(t,
		map[string]any{"module": "auth", "created_at": now},
		map[string]any{"module": "auth", "created_at": now},
		map[string]any{"module": "platform", "created_at": now},
	)

	from := now.Add(-1 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
		WithQuery("from", from).
		WithQuery("to", to).
		WithQuery("module", "auth").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	arr := resp.JSON().Array()
	arr.Length().IsEqual(2)
	arr.Value(0).Object().Value("module").String().IsEqual("auth")
	arr.Value(1).Object().Value("module").String().IsEqual("auth")
}

func TestGetActionLogs_FilterByTags(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionActionLogRead)

	now := time.Now()
	audit.GivenActionLogs(t,
		map[string]any{"module": "auth", "tags": []string{"tender", "procurement"}, "created_at": now},
		map[string]any{"module": "auth", "tags": []string{"tender"}, "created_at": now},
		map[string]any{"module": "auth", "tags": []string{"user"}, "created_at": now},
		map[string]any{"module": "auth", "created_at": now}, // no tags
	)

	from := now.Add(-1 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)

	// WHEN - filter by single tag "tender"
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
		WithQuery("from", from).
		WithQuery("to", to).
		WithQuery("tags", "tender").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN - only logs containing "tender" tag are returned
	resp.Status(http.StatusOK)
	arr := resp.JSON().Array()
	arr.Length().IsEqual(2)
}

func TestGetActionLogs_FilterByGroupKey(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionActionLogRead)

	groupKey1 := "tender:" + uuid.NewString()
	groupKey2 := "tender:" + uuid.NewString()

	now := time.Now()
	audit.GivenActionLogs(t,
		map[string]any{"module": "auth", "group_key": groupKey1, "created_at": now},
		map[string]any{"module": "auth", "group_key": groupKey1, "created_at": now},
		map[string]any{"module": "auth", "group_key": groupKey2, "created_at": now},
		map[string]any{"module": "auth", "created_at": now}, // no group_key
	)

	from := now.Add(-1 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)

	// WHEN - filter by group_key
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
		WithQuery("from", from).
		WithQuery("to", to).
		WithQuery("group_key", groupKey1).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN - only logs with matching group_key are returned
	resp.Status(http.StatusOK)
	arr := resp.JSON().Array()
	arr.Length().IsEqual(2)
}

func TestGetActionLogs_CursorPagination(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionActionLogRead)

	now := time.Now()
	ids := audit.GivenActionLogs(t,
		map[string]any{"operation_id": "op-1", "created_at": now.Add(-3 * time.Second)},
		map[string]any{"operation_id": "op-2", "created_at": now.Add(-2 * time.Second)},
		map[string]any{"operation_id": "op-3", "created_at": now.Add(-1 * time.Second)},
	)

	from := now.Add(-1 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)

	// WHEN - first page (limit=2, no cursor)
	resp1 := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
		WithQuery("from", from).
		WithQuery("to", to).
		WithQuery("limit", 2).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN - first page returns 2 items (newest first: op-3, op-2)
	resp1.Status(http.StatusOK)
	arr1 := resp1.JSON().Array()
	arr1.Length().IsEqual(2)
	arr1.Value(0).Object().Value("operation_id").String().IsEqual("op-3")
	arr1.Value(1).Object().Value("operation_id").String().IsEqual("op-2")

	// WHEN - second page (cursor = id of second item from first page, which is ids[1])
	resp2 := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
		WithQuery("from", from).
		WithQuery("to", to).
		WithQuery("limit", 2).
		WithQuery("cursor", ids[1]).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN - second page returns remaining 1 item (op-1)
	resp2.Status(http.StatusOK)
	arr2 := resp2.JSON().Array()
	arr2.Length().IsEqual(1)
	arr2.Value(0).Object().Value("operation_id").String().IsEqual("op-1")
}

func TestGetActionLogs_MissingTimeRange(t *testing.T) {
	tests := []struct {
		name       string
		from       *string
		to         *string
		wantStatus int
	}{
		{
			name:       "missing from",
			from:       nil,
			to:         lo.ToPtr(time.Now().Format(time.RFC3339)),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing to",
			from:       lo.ToPtr(time.Now().Format(time.RFC3339)),
			to:         nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing both",
			from:       nil,
			to:         nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := auth.GivenAuthToken(t, portalaudit.PermissionActionLogRead)

			// WHEN
			req := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
				WithHeader("Authorization", "Bearer "+token)

			if tc.from != nil {
				req = req.WithQuery("from", *tc.from)
			}
			if tc.to != nil {
				req = req.WithQuery("to", *tc.to)
			}
			resp := req.Expect()

			// THEN
			resp.Status(tc.wantStatus)
		})
	}
}

func TestGetActionLogs_AuthFailures(t *testing.T) {
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
				return auth.GivenAuthToken(t)
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

			now := time.Now()
			from := now.Add(-1 * time.Hour).Format(time.RFC3339)
			to := now.Add(1 * time.Hour).Format(time.RFC3339)

			// WHEN
			req := trigger.UserAction(t).GET("/api/v1/audit/get-action-logs").
				WithQuery("from", from).
				WithQuery("to", to)
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
