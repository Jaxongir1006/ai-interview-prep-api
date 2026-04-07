//go:build system

package trigger

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"

	"github.com/rise-and-shine/pkg/taskmill/enqueuer"
	"github.com/uptrace/bun"
)

const (
	asyncTaskTimeout  = 10 * time.Second
	asyncPollInterval = 100 * time.Millisecond
	taskmillSchema    = "taskmill"
	taskQueueTable    = "task_queue"
)

// AsyncTask enqueues a task and waits for it to complete.
// Returns the task ID on success. Fails the test if the task doesn't complete
// within timeout or moves to DLQ.
// Works for both ephemeral and non-ephemeral tasks.
func AsyncTask(t *testing.T, queueName, operationID string, payload any) int64 {
	t.Helper()

	db := database.GetTestDB(t)

	eq, err := enqueuer.New(queueName)
	if err != nil {
		t.Fatalf("AsyncTask: failed to create enqueuer: %v", err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), asyncTaskTimeout)
	defer cancel()

	taskID, err := eq.Enqueue(ctx, db, operationID, payload)
	if err != nil {
		t.Fatalf("AsyncTask: failed to enqueue task: %v", err)
	}

	waitForTaskCompletion(t, db, taskID, asyncTaskTimeout)

	return taskID
}

type taskStatus int

const (
	taskStatusPending taskStatus = iota
	taskStatusCompleted
	taskStatusDLQ
)

// waitForTaskCompletion polls the task_queue table until the task is removed (success)
// or moved to DLQ (failure).
func waitForTaskCompletion(t *testing.T, db *bun.DB, taskID int64, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		status, err := getTaskStatus(db, taskID)
		if err != nil {
			t.Fatalf("waitForTaskCompletion: failed to check task status: %v", err)
		}

		switch status {
		case taskStatusCompleted:
			return
		case taskStatusDLQ:
			t.Fatalf("waitForTaskCompletion: task %d moved to DLQ", taskID)
		case taskStatusPending:
			time.Sleep(asyncPollInterval)
		}
	}

	t.Fatalf("waitForTaskCompletion: task %d did not complete within %v", taskID, timeout)
}

func getTaskStatus(db *bun.DB, taskID int64) (taskStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dlqAt *time.Time
	err := db.NewSelect().
		TableExpr(taskmillSchema+"."+taskQueueTable).
		Column("dlq_at").
		Where("id = ?", taskID).
		Scan(ctx, &dlqAt)

	if err == sql.ErrNoRows {
		return taskStatusCompleted, nil
	}
	if err != nil {
		return 0, err
	}
	if dlqAt != nil {
		return taskStatusDLQ, nil
	}
	return taskStatusPending, nil
}
