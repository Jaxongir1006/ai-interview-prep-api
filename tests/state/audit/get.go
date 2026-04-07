package audit

import (
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/actionlog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
)

// GetActionLogByID retrieves an action log by ID.
// Fails the test if the action log is not found.
func GetActionLogByID(t *testing.T, id int64) *actionlog.ActionLog {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewActionLogRepo(db, db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	log, err := repo.Get(ctx, actionlog.Filter{ID: &id})
	if err != nil {
		t.Fatalf("GetActionLogByID: failed to get action log %d: %v", id, err)
	}

	return log
}

// GetActionLogs retrieves all action logs ordered by id DESC.
func GetActionLogs(t *testing.T) []actionlog.ActionLog {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewActionLogRepo(db, db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	logs, err := repo.List(ctx, actionlog.Filter{})
	if err != nil {
		t.Fatalf("GetActionLogs: failed to list action logs: %v", err)
	}

	return logs
}

// ActionLogCount returns the total count of action logs.
func ActionLogCount(t *testing.T) int {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewActionLogRepo(db, db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := repo.Count(ctx, actionlog.Filter{})
	if err != nil {
		t.Fatalf("ActionLogCount: failed to count action logs: %v", err)
	}

	return count
}

// StatusChangeLogCount returns the total count of status change logs.
func StatusChangeLogCount(t *testing.T) int {
	t.Helper()

	db := database.GetTestDB(t)

	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := db.NewSelect().
		TableExpr("audit.status_change_logs").
		Count(ctx)
	if err != nil {
		t.Fatalf("StatusChangeLogCount: failed to count status change logs: %v", err)
	}

	return count
}
