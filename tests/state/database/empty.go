package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/uptrace/bun"
)

// appSchemas returns the application schemas to clean.
// Add new schemas here as the application grows.
func appSchemas() []string {
	return []string{"auth", "audit", "candidate", "analytics", "taskmill", "alert", "filevault", "interview"}
}

// Empty truncates all application tables, providing a clean state for tests.
func Empty(t *testing.T) {
	t.Helper()

	EmptyRedis(t)

	db := GetTestDB(t)

	ctx, cancel := QueryContext()
	defer cancel()

	tables, err := discoverTables(ctx, db)
	if err != nil {
		t.Fatalf("failed to discover tables: %v", err)
	}

	if len(tables) == 0 {
		return
	}

	query := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))

	if _, err = db.ExecContext(ctx, query); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func EmptyRedis(t *testing.T) {
	t.Helper()

	client := GetTestRedis(t)

	ctx, cancel := QueryContext()
	defer cancel()

	err := client.FlushDB(ctx).Err()
	if err != nil {
		t.Fatalf("failed to flush Redis: %v", err)
	}
}

func discoverTables(ctx context.Context, db *bun.DB) ([]string, error) {
	var tables []string

	err := db.NewSelect().
		ColumnExpr("table_schema || '.' || table_name").
		TableExpr("information_schema.tables").
		Where("table_schema IN (?)", bun.In(appSchemas())).
		Where("table_type = ?", "BASE TABLE").
		Scan(ctx, &tables)
	if err != nil {
		return nil, err
	}

	return tables, nil
}
