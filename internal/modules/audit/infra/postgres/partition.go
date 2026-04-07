package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/code19m/errx"
	"github.com/uptrace/bun"
)

//nolint:gochecknoglobals // package-level cache for ensured partitions
var ensuredPartitions sync.Map

var partitionTables = []string{"action_logs", "status_change_logs"} //nolint:gochecknoglobals // static list

// createWithPartitions ensures monthly partitions exist, then delegates to the inner create.
func createWithPartitions[E any](
	ctx context.Context,
	db *bun.DB,
	createFn func(context.Context, *E) (*E, error),
	entity *E,
) (*E, error) {
	err := ensurePartitions(ctx, db, time.Now())
	if err != nil {
		return nil, errx.Wrap(err)
	}
	return createFn(ctx, entity)
}

// ensurePartitions creates monthly partitions for the given time if they don't already exist.
// Safe to call concurrently — uses CREATE TABLE IF NOT EXISTS and an in-memory cache.
func ensurePartitions(ctx context.Context, db *bun.DB, t time.Time) error {
	year, month, _ := t.Date()
	key := fmt.Sprintf("y%dm%02d", year, month)

	if _, ok := ensuredPartitions.Load(key); ok {
		return nil
	}

	rangeStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	rangeEnd := rangeStart.AddDate(0, 1, 0)

	for _, table := range partitionTables {
		_, err := db.ExecContext(ctx, fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s.%s_%s PARTITION OF %s.%s FOR VALUES FROM ('%s') TO ('%s')`,
			schemaName, table, key,
			schemaName, table,
			rangeStart.Format("2006-01-02"),
			rangeEnd.Format("2006-01-02"),
		))
		if err != nil {
			return errx.Wrap(err)
		}
	}

	ensuredPartitions.Store(key, struct{}{})
	return nil
}
