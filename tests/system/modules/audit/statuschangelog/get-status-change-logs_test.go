//go:build system

package statuschangelog_test

import (
	"net/http"
	"testing"
	"time"

	portalaudit "go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/tests/state/audit"
	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	"go-enterprise-blueprint/tests/system/trigger"
)

func TestGetStatusChangeLogs_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionStatusChangeLogRead)

	now := time.Now().UTC().Truncate(time.Second)
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	actionLogIDs := audit.GivenActionLogs(t,
		map[string]any{"created_at": now},
	)

	audit.GivenStatusChangeLogs(t,
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-1",
			"status":        "approved",
			"created_at":    now,
		},
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-2",
			"status":        "rejected",
			"created_at":    now,
		},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-status-change-logs").
		WithQuery("from", from.Format(time.RFC3339)).
		WithQuery("to", to.Format(time.RFC3339)).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	arr := resp.JSON().Array()
	arr.Length().IsEqual(2)
	arr.Value(0).Object().Value("entity_type").String().IsEqual("order")
	arr.Value(0).Object().Value("status").String().NotEmpty()
	arr.Value(0).Object().Value("action_log_id").Number().Gt(0)
	arr.Value(0).Object().Value("trace_id").String().NotEmpty()
	arr.Value(0).Object().Value("created_at").String().NotEmpty()
}

func TestGetStatusChangeLogs_FilterByEntityType(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionStatusChangeLogRead)

	now := time.Now().UTC().Truncate(time.Second)
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	actionLogIDs := audit.GivenActionLogs(t,
		map[string]any{"created_at": now},
	)

	audit.GivenStatusChangeLogs(t,
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-1",
			"status":        "approved",
			"created_at":    now,
		},
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "user",
			"entity_id":     "user-1",
			"status":        "active",
			"created_at":    now,
		},
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-2",
			"status":        "rejected",
			"created_at":    now,
		},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-status-change-logs").
		WithQuery("from", from.Format(time.RFC3339)).
		WithQuery("to", to.Format(time.RFC3339)).
		WithQuery("entity_type", "order").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	arr := resp.JSON().Array()
	arr.Length().IsEqual(2)
	arr.Value(0).Object().Value("entity_type").String().IsEqual("order")
	arr.Value(1).Object().Value("entity_type").String().IsEqual("order")
}

func TestGetStatusChangeLogs_CursorPagination(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionStatusChangeLogRead)

	now := time.Now().UTC().Truncate(time.Second)
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	actionLogIDs := audit.GivenActionLogs(t,
		map[string]any{"created_at": now},
	)

	logIDs := audit.GivenStatusChangeLogs(t,
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-1",
			"status":        "created",
			"created_at":    now,
		},
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-2",
			"status":        "approved",
			"created_at":    now,
		},
		map[string]any{
			"action_log_id": actionLogIDs[0],
			"entity_type":   "order",
			"entity_id":     "order-3",
			"status":        "shipped",
			"created_at":    now,
		},
	)

	// WHEN - first page with limit 2
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-status-change-logs").
		WithQuery("from", from.Format(time.RFC3339)).
		WithQuery("to", to.Format(time.RFC3339)).
		WithQuery("limit", 2).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Array().Length().IsEqual(2)

	// WHEN - second page using cursor
	resp2 := trigger.UserAction(t).GET("/api/v1/audit/get-status-change-logs").
		WithQuery("from", from.Format(time.RFC3339)).
		WithQuery("to", to.Format(time.RFC3339)).
		WithQuery("cursor", logIDs[1]).
		WithQuery("limit", 2).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp2.Status(http.StatusOK)
	resp2.JSON().Array().Length().IsEqual(1)
}

func TestGetStatusChangeLogs_MissingTimeRange(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t, portalaudit.PermissionStatusChangeLogRead)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/audit/get-status-change-logs").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusBadRequest)
	resp.JSON().Object().Value("error").Object().
		Value("code").String().IsEqual("VALIDATION_FAILED")
}

func TestGetStatusChangeLogs_AuthFailures(t *testing.T) {
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

			// WHEN
			req := trigger.UserAction(t).GET("/api/v1/audit/get-status-change-logs").
				WithQuery("from", "2024-01-01T00:00:00Z").
				WithQuery("to", "2024-12-31T23:59:59Z")
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
