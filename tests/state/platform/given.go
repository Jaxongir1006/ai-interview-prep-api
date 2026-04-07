package platform

import (
	"testing"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/anymap"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/uptrace/bun"
)

// taskQueueRecord represents a record in taskmill.task_queue table.
type taskQueueRecord struct {
	bun.BaseModel `bun:"table:taskmill.task_queue,alias:tq"`

	ID             int64          `bun:"id,pk,autoincrement"`
	QueueName      string         `bun:"queue_name,notnull"`
	TaskGroupID    *string        `bun:"task_group_id"`
	OperationID    string         `bun:"operation_id,notnull"`
	Meta           map[string]any `bun:"meta,type:jsonb"`
	Payload        map[string]any `bun:"payload,type:jsonb,notnull"`
	ScheduledAt    time.Time      `bun:"scheduled_at,notnull,default:current_timestamp"`
	VisibleAt      time.Time      `bun:"visible_at,notnull,default:current_timestamp"`
	ExpiresAt      *time.Time     `bun:"expires_at"`
	Priority       int            `bun:"priority,notnull,default:0"`
	Attempts       int            `bun:"attempts,notnull,default:0"`
	MaxAttempts    int            `bun:"max_attempts,notnull,default:3"`
	IdempotencyKey string         `bun:"idempotency_key,notnull"`
	CreatedAt      time.Time      `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time      `bun:"updated_at,notnull,default:current_timestamp"`
	DLQAt          *time.Time     `bun:"dlq_at"`
	DLQReason      map[string]any `bun:"dlq_reason,type:jsonb"`
	Ephemeral      bool           `bun:"ephemeral,notnull,default:false"`
}

// taskResultRecord represents a record in taskmill.task_results table.
type taskResultRecord struct {
	bun.BaseModel `bun:"table:taskmill.task_results,alias:tr"`

	ID             int64          `bun:"id,pk"`
	QueueName      string         `bun:"queue_name,notnull"`
	TaskGroupID    *string        `bun:"task_group_id"`
	OperationID    string         `bun:"operation_id,notnull"`
	Meta           map[string]any `bun:"meta,type:jsonb"`
	Payload        map[string]any `bun:"payload,type:jsonb,notnull"`
	Priority       int            `bun:"priority,notnull"`
	Attempts       int            `bun:"attempts,notnull"`
	MaxAttempts    int            `bun:"max_attempts,notnull"`
	IdempotencyKey string         `bun:"idempotency_key,notnull"`
	ScheduledAt    time.Time      `bun:"scheduled_at,notnull"`
	CreatedAt      time.Time      `bun:"created_at,notnull"`
	CompletedAt    time.Time      `bun:"completed_at,notnull,default:current_timestamp"`
}

// taskScheduleRecord represents a record in taskmill.task_schedules table.
type taskScheduleRecord struct {
	bun.BaseModel `bun:"table:taskmill.task_schedules,alias:ts"`

	ID            int64      `bun:"id,pk,autoincrement"`
	OperationID   string     `bun:"operation_id,notnull,unique"`
	QueueName     string     `bun:"queue_name,notnull"`
	CronPattern   string     `bun:"cron_pattern,notnull"`
	NextRunAt     time.Time  `bun:"next_run_at,notnull"`
	LastRunAt     *time.Time `bun:"last_run_at"`
	LastRunStatus *string    `bun:"last_run_status"`
	LastError     *string    `bun:"last_error"`
	RunCount      int64      `bun:"run_count,notnull,default:0"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
}

//nolint:gochecknoglobals // static validation maps for test state
var (
	validQueuedTaskKeys = map[string]bool{
		"queue_name": true, "task_group_id": true, "operation_id": true, "meta": true, "payload": true,
		"scheduled_at": true, "visible_at": true, "expires_at": true, "priority": true, "attempts": true,
		"max_attempts": true, "idempotency_key": true,
	}

	validDLQTaskKeys = map[string]bool{
		"queue_name": true, "task_group_id": true, "operation_id": true, "meta": true, "payload": true,
		"dlq_reason": true, "priority": true, "attempts": true, "max_attempts": true, "idempotency_key": true,
	}

	validTaskResultKeys = map[string]bool{
		"queue_name": true, "task_group_id": true, "operation_id": true, "meta": true, "payload": true,
		"priority": true, "attempts": true, "max_attempts": true, "idempotency_key": true,
		"scheduled_at": true, "created_at": true, "completed_at": true,
	}

	validScheduleKeys = map[string]bool{
		"operation_id": true, "queue_name": true, "cron_pattern": true, "next_run_at": true,
		"last_run_at": true, "last_run_status": true, "last_error": true, "run_count": true,
	}
)

// GivenQueuedTasks creates task records in taskmill.task_queue.
// Valid keys: queue_name (required), operation_id (default: "test-op"),
// payload (default: {}), priority (default: 0), max_attempts (default: 3),
// scheduled_at (default: now), idempotency_key (default: uuid).
func GivenQueuedTasks(t *testing.T, data ...map[string]any) []int64 {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenQueuedTasks: at least one task data map is required")
	}

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]int64, 0, len(data))
	now := time.Now()

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenQueuedTasks", validQueuedTaskKeys, d)

		queueName := anymap.String(d, "queue_name", "")
		if queueName == "" {
			t.Fatalf("GivenQueuedTasks[%d]: queue_name is required", i)
		}

		payload := getMapOrDefault(d, "payload", map[string]any{})

		record := &taskQueueRecord{
			QueueName:      queueName,
			TaskGroupID:    anymap.StringPtr(d, "task_group_id", nil),
			OperationID:    anymap.String(d, "operation_id", "test-op"),
			Meta:           getMapOrNil(d, "meta"),
			Payload:        payload,
			ScheduledAt:    anymap.Time(d, "scheduled_at", now),
			VisibleAt:      anymap.Time(d, "visible_at", now),
			ExpiresAt:      anymap.TimePtr(d, "expires_at", nil),
			Priority:       getIntOrDefault(d, "priority", 0),
			Attempts:       getIntOrDefault(d, "attempts", 0),
			MaxAttempts:    getIntOrDefault(d, "max_attempts", 3),
			IdempotencyKey: anymap.String(d, "idempotency_key", uuid.NewString()),
			CreatedAt:      now,
			UpdatedAt:      now,
			Ephemeral:      false,
		}

		var id int64
		err := db.NewInsert().
			Model(record).
			Returning("id").
			Scan(ctx, &id)

		if err != nil {
			t.Fatalf("GivenQueuedTasks[%d]: failed to insert task: %v", i, err)
		}

		ids = append(ids, id)
	}

	return ids
}

// GivenDLQTasks creates DLQ task records (tasks with dlq_at set).
// Valid keys: queue_name (required), operation_id (default: "test-op"),
// payload (default: {}), dlq_reason (default: {"error": "test failure"}),
// priority, attempts, max_attempts, idempotency_key.
func GivenDLQTasks(t *testing.T, data ...map[string]any) []int64 {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenDLQTasks: at least one task data map is required")
	}

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]int64, 0, len(data))
	now := time.Now()

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenDLQTasks", validDLQTaskKeys, d)

		queueName := anymap.String(d, "queue_name", "")
		if queueName == "" {
			t.Fatalf("GivenDLQTasks[%d]: queue_name is required", i)
		}

		payload := getMapOrDefault(d, "payload", map[string]any{})
		dlqReason := getMapOrDefault(d, "dlq_reason", map[string]any{"error": "test failure"})

		record := &taskQueueRecord{
			QueueName:      queueName,
			TaskGroupID:    anymap.StringPtr(d, "task_group_id", nil),
			OperationID:    anymap.String(d, "operation_id", "test-op"),
			Meta:           getMapOrNil(d, "meta"),
			Payload:        payload,
			ScheduledAt:    now,
			VisibleAt:      now,
			Priority:       getIntOrDefault(d, "priority", 0),
			Attempts:       getIntOrDefault(d, "attempts", 3),
			MaxAttempts:    getIntOrDefault(d, "max_attempts", 3),
			IdempotencyKey: anymap.String(d, "idempotency_key", uuid.NewString()),
			CreatedAt:      now,
			UpdatedAt:      now,
			DLQAt:          lo.ToPtr(now),
			DLQReason:      dlqReason,
			Ephemeral:      false,
		}

		var id int64
		err := db.NewInsert().
			Model(record).
			Returning("id").
			Scan(ctx, &id)

		if err != nil {
			t.Fatalf("GivenDLQTasks[%d]: failed to insert DLQ task: %v", i, err)
		}

		ids = append(ids, id)
	}

	return ids
}

// GivenTaskResults creates records in taskmill.task_results.
// Valid keys: queue_name (required), operation_id (default: "test-op"),
// payload (default: {}), attempts (default: 1), max_attempts (default: 3),
// completed_at (default: now), idempotency_key (default: uuid).
func GivenTaskResults(t *testing.T, data ...map[string]any) []int64 {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenTaskResults: at least one result data map is required")
	}

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]int64, 0, len(data))
	now := time.Now()

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenTaskResults", validTaskResultKeys, d)

		queueName := anymap.String(d, "queue_name", "")
		if queueName == "" {
			t.Fatalf("GivenTaskResults[%d]: queue_name is required", i)
		}

		payload := getMapOrDefault(d, "payload", map[string]any{})

		// Generate a unique ID for task_results (not auto-increment)
		id := generateTaskResultID()

		record := &taskResultRecord{
			ID:             id,
			QueueName:      queueName,
			TaskGroupID:    anymap.StringPtr(d, "task_group_id", nil),
			OperationID:    anymap.String(d, "operation_id", "test-op"),
			Meta:           getMapOrNil(d, "meta"),
			Payload:        payload,
			Priority:       getIntOrDefault(d, "priority", 0),
			Attempts:       getIntOrDefault(d, "attempts", 1),
			MaxAttempts:    getIntOrDefault(d, "max_attempts", 3),
			IdempotencyKey: anymap.String(d, "idempotency_key", uuid.NewString()),
			ScheduledAt:    anymap.Time(d, "scheduled_at", now.Add(-1*time.Hour)),
			CreatedAt:      anymap.Time(d, "created_at", now.Add(-1*time.Hour)),
			CompletedAt:    anymap.Time(d, "completed_at", now),
		}

		_, err := db.NewInsert().
			Model(record).
			Exec(ctx)

		if err != nil {
			t.Fatalf("GivenTaskResults[%d]: failed to insert task result: %v", i, err)
		}

		ids = append(ids, id)
	}

	return ids
}

// GivenSchedules creates records in taskmill.task_schedules.
// Valid keys: operation_id (required), queue_name (required),
// cron_pattern (default: "0 * * * *"), next_run_at (default: now+1h),
// run_count (default: 0).
func GivenSchedules(t *testing.T, data ...map[string]any) []int64 {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenSchedules: at least one schedule data map is required")
	}

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	ids := make([]int64, 0, len(data))
	now := time.Now()

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenSchedules", validScheduleKeys, d)

		operationID := anymap.String(d, "operation_id", "")
		queueName := anymap.String(d, "queue_name", "")

		if operationID == "" || queueName == "" {
			t.Fatalf("GivenSchedules[%d]: operation_id and queue_name are required", i)
		}

		record := &taskScheduleRecord{
			OperationID:   operationID,
			QueueName:     queueName,
			CronPattern:   anymap.String(d, "cron_pattern", "0 * * * *"),
			NextRunAt:     anymap.Time(d, "next_run_at", now.Add(1*time.Hour)),
			LastRunAt:     anymap.TimePtr(d, "last_run_at", nil),
			LastRunStatus: anymap.StringPtr(d, "last_run_status", nil),
			LastError:     anymap.StringPtr(d, "last_error", nil),
			RunCount:      int64(getIntOrDefault(d, "run_count", 0)),
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		var id int64
		err := db.NewInsert().
			Model(record).
			Returning("id").
			Scan(ctx, &id)

		if err != nil {
			t.Fatalf("GivenSchedules[%d]: failed to insert schedule: %v", i, err)
		}

		ids = append(ids, id)
	}

	return ids
}

// generateTaskResultID generates a unique ID for task results.
// taskmill uses snowflake IDs, but for tests we can use a simpler approach.
func generateTaskResultID() int64 {
	return time.Now().UnixNano() / 1000000 // milliseconds
}

// GetTaskQueueCount returns the count of tasks in taskmill.task_queue.
func GetTaskQueueCount(t *testing.T, queueName string) int {
	t.Helper()

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := db.NewSelect().
		TableExpr("taskmill.task_queue").
		Where("queue_name = ?", queueName).
		Where("dlq_at IS NULL").
		Count(ctx)

	if err != nil {
		t.Fatalf("GetTaskQueueCount: failed to count tasks: %v", err)
	}

	return count
}

// GetDLQCount returns the count of DLQ tasks in taskmill.task_queue.
func GetDLQCount(t *testing.T, queueName string) int {
	t.Helper()

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := db.NewSelect().
		TableExpr("taskmill.task_queue").
		Where("queue_name = ?", queueName).
		Where("dlq_at IS NOT NULL").
		Count(ctx)

	if err != nil {
		t.Fatalf("GetDLQCount: failed to count DLQ tasks: %v", err)
	}

	return count
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

// getMapOrNil extracts a map[string]any from data or returns nil.
func getMapOrNil(data map[string]any, key string) map[string]any {
	if v, hasKey := data[key]; hasKey {
		if m, isMap := v.(map[string]any); isMap {
			return m
		}
	}
	return nil
}

// getIntOrDefault extracts an int from data or returns default.
func getIntOrDefault(data map[string]any, key string, defaultVal int) int {
	if v, hasKey := data[key]; hasKey {
		return cast.ToInt(v)
	}
	return defaultVal
}
