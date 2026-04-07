package platform

import (
	"sync"
	"testing"
	"time"

	"github.com/jaxongir1006/hire-ready-api/pkg/anymap"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

//nolint:gochecknoglobals // one-time schema init for tests
var alertSchemaOnce sync.Once

// ensureAlertSchema creates the alert schema and errors table if they don't exist.
// The alert provider normally creates this at startup, but the test config uses "noop" provider.
func ensureAlertSchema(t *testing.T) {
	t.Helper()

	alertSchemaOnce.Do(func() {
		db := database.GetTestDB(t)
		ctx, cancel := database.QueryContext()
		defer cancel()

		_, err := db.ExecContext(ctx, `
			CREATE SCHEMA IF NOT EXISTS alert;
			CREATE TABLE IF NOT EXISTS alert.errors (
				id UUID PRIMARY KEY,
				code TEXT NOT NULL,
				message TEXT NOT NULL,
				details JSONB NOT NULL,
				service TEXT NOT NULL,
				operation TEXT NOT NULL,
				created_at TIMESTAMPTZ NOT NULL,
				alerted BOOLEAN NOT NULL
			);
			CREATE INDEX IF NOT EXISTS idx_errors_service_operation_alerted ON alert.errors (service, operation, alerted);
			CREATE INDEX IF NOT EXISTS idx_errors_created_at ON alert.errors (created_at);
		`)
		if err != nil {
			t.Fatalf("ensureAlertSchema: failed to create alert schema/table: %v", err)
		}
	})
}

// errorRecord represents a record in alert.errors table.
type errorRecord struct {
	bun.BaseModel `bun:"table:alert.errors,alias:ae"`

	ID        string            `bun:"id,pk"`
	Code      string            `bun:"code,notnull"`
	Message   string            `bun:"message,notnull"`
	Details   map[string]string `bun:"details,type:jsonb,notnull"`
	Service   string            `bun:"service,notnull"`
	Operation string            `bun:"operation,notnull"`
	CreatedAt time.Time         `bun:"created_at,notnull"`
	Alerted   bool              `bun:"alerted,notnull"`
}

//nolint:gochecknoglobals // static validation map for test state
var validErrorKeys = map[string]bool{
	"code": true, "service": true, "operation": true, "message": true,
	"details": true, "alerted": true, "created_at": true,
}

// GivenErrors creates error records in alert.errors.
// Valid keys: code (default: "INTERNAL_ERROR"), service (default: "blueprint"),
// operation (default: "test-op"), message (default: "test error"),
// details (default: {}), alerted (default: false), created_at (default: now).
// Returns []string of generated UUIDs.
func GivenErrors(t *testing.T, data ...map[string]any) []string {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenErrors: at least one error data map is required")
	}

	ensureAlertSchema(t)

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]string, 0, len(data))
	now := time.Now()

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenErrors", validErrorKeys, d)

		id := uuid.NewString()

		record := &errorRecord{
			ID:        id,
			Code:      anymap.String(d, "code", "INTERNAL_ERROR"),
			Message:   anymap.String(d, "message", "test error"),
			Details:   getStringMapOrDefault(d, "details", map[string]string{}),
			Service:   anymap.String(d, "service", "blueprint"),
			Operation: anymap.String(d, "operation", "test-op"),
			CreatedAt: anymap.Time(d, "created_at", now),
			Alerted:   anymap.Bool(d, "alerted", false),
		}

		_, err := db.NewInsert().
			Model(record).
			Exec(ctx)

		if err != nil {
			t.Fatalf("GivenErrors[%d]: failed to insert error: %v", i, err)
		}

		ids = append(ids, id)
	}

	return ids
}

// GetErrorCount returns the count of all errors in alert.errors.
func GetErrorCount(t *testing.T) int {
	t.Helper()

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := db.NewSelect().
		TableExpr("alert.errors").
		Count(ctx)

	if err != nil {
		t.Fatalf("GetErrorCount: failed to count errors: %v", err)
	}

	return count
}

// ErrorExists checks if an error record exists in alert.errors by ID.
func ErrorExists(t *testing.T, id string) bool {
	t.Helper()

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	exists, err := db.NewSelect().
		TableExpr("alert.errors").
		Where("id = ?", id).
		Exists(ctx)

	if err != nil {
		t.Fatalf("ErrorExists: failed to check error existence: %v", err)
	}

	return exists
}

// getStringMapOrDefault extracts a map[string]string from data or returns default.
func getStringMapOrDefault(data map[string]any, key string, defaultVal map[string]string) map[string]string {
	if v, hasKey := data[key]; hasKey {
		if m, isMap := v.(map[string]string); isMap {
			return m
		}
	}
	return defaultVal
}
