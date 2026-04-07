package audit

import (
	"testing"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/actionlog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/statuschangelog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/anymap"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// GivenActionLogs creates action log records in the database for test setup.
// Each map in data represents an action log with the following valid keys:
//   - user_id: *string (default: nil)
//   - module: string (default: "auth")
//   - operation_id: string (default: "test-op")
//   - request_payload: map[string]any (default: {})
//   - tags: []string (default: nil)
//   - group_key: *string (default: nil)
//   - ip_address: string (default: "127.0.0.1")
//   - user_agent: string (default: "test-agent")
//   - trace_id: string (default: "test-trace")
//   - created_at: time.Time (default: now)
//
// Returns the created action log IDs.
func GivenActionLogs(t *testing.T, data ...map[string]any) []int64 {
	t.Helper()

	if len(data) == 0 {
		data = []map[string]any{{}}
	}

	db := database.GetTestDB(t)
	repo := postgres.NewActionLogRepo(db, db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]int64, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenActionLogs", validActionLogKeys, d)

		entity := &actionlog.ActionLog{
			UserID:         anymap.StringPtr(d, "user_id", nil),
			Module:         anymap.String(d, "module", "auth"),
			OperationID:    anymap.String(d, "operation_id", "test-op"),
			RequestPayload: getMapOrDefault(d, "request_payload", map[string]any{}),
			Tags:           anymap.StringSlice(d, "tags", nil),
			GroupKey:       anymap.StringPtr(d, "group_key", nil),
			IPAddress:      anymap.String(d, "ip_address", "127.0.0.1"),
			UserAgent:      anymap.String(d, "user_agent", "test-agent"),
			TraceID:        anymap.String(d, "trace_id", "test-trace"),
			CreatedAt:      anymap.Time(d, "created_at", time.Now()),
		}

		created, err := repo.Create(ctx, entity)
		if err != nil {
			t.Fatalf("GivenActionLogs[%d]: failed to insert action log: %v", i, err)
		}

		ids = append(ids, created.ID)
	}

	return ids
}

// GivenStatusChangeLogs creates status change log records in the database for test setup.
// Each map in data represents a status change log with the following valid keys:
//   - action_log_id: int64 (required)
//   - entity_type: string (default: "user")
//   - entity_id: string (default: generated UUID)
//   - status: string (default: "active")
//   - trace_id: string (default: "test-trace")
//   - created_at: time.Time (default: now)
//
// Returns the created status change log IDs.
func GivenStatusChangeLogs(t *testing.T, data ...map[string]any) []int64 {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenStatusChangeLogs: at least one data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewStatusChangeLogRepo(db, db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]int64, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenStatusChangeLogs", validStatusChangeLogKeys, d)

		actionLogID := cast.ToInt64(d["action_log_id"])
		if actionLogID == 0 {
			t.Fatalf("GivenStatusChangeLogs[%d]: action_log_id is required", i)
		}

		entity := &statuschangelog.StatusChangeLog{
			ActionLogID: actionLogID,
			EntityType:  anymap.String(d, "entity_type", "user"),
			EntityID:    anymap.String(d, "entity_id", uuid.NewString()),
			Status:      anymap.String(d, "status", "active"),
			TraceID:     anymap.String(d, "trace_id", "test-trace"),
			CreatedAt:   anymap.Time(d, "created_at", time.Now()),
		}

		created, err := repo.Create(ctx, entity)
		if err != nil {
			t.Fatalf("GivenStatusChangeLogs[%d]: failed to insert status change log: %v", i, err)
		}

		ids = append(ids, created.ID)
	}

	return ids
}

// getMapOrDefault extracts a map[string]any from data or returns default.
func getMapOrDefault(data map[string]any, key string, defaultVal map[string]any) map[string]any {
	if v, hasKey := data[key]; hasKey {
		if m, isMap := v.(map[string]any); isMap {
			return m
		}
	}
	return defaultVal
}
