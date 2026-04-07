package database

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/pg"
	"github.com/uptrace/bun"
)

//nolint:gochecknoglobals // lazy singleton for test database connection
var (
	dbOnce sync.Once
	dbInst *bun.DB
)

// config is a minimal configuration that reads from test.yaml.
type config struct {
	Postgres pg.Config `yaml:"postgres" validate:"required"`
}

// GetTestDB returns a lazily initialized singleton *bun.DB connected to the test database.
// Fails if database name doesn't start with "test_" to prevent accidental data loss.
func GetTestDB(t *testing.T) *bun.DB {
	t.Helper()

	dbOnce.Do(func() {
		db, err := initDB()
		if err != nil {
			t.Fatalf("failed to initialize test database: %v", err)
		}
		dbInst = db
	})

	return dbInst
}

// QueryContext returns a context with 10 second timeout for database operations.
func QueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func initDB() (db *bun.DB, err error) { //nolint:nonamedreturns // required to capture deferred error
	originalWd, err := os.Getwd()
	if err != nil {
		return nil, errx.Newf("get working directory: %v", err)
	}

	root, err := projectRoot()
	if err != nil {
		return nil, errx.Newf("find project root: %v", err)
	}

	if err = os.Chdir(root); err != nil {
		return nil, errx.Newf("change to project root: %v", err)
	}
	defer func() { err = os.Chdir(originalWd) }()

	cfg := cfgloader.MustLoad[config](cfgloader.WithSilent())

	if !strings.HasPrefix(cfg.Postgres.Database, "test_") {
		return nil, errx.Newf(
			"SAFETY: refusing to use database %q - name must start with 'test_'",
			cfg.Postgres.Database,
		)
	}

	return pg.NewBunDB(cfg.Postgres)
}

func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errx.Newf("get working directory: %v", err)
	}

	for {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errx.Newf("could not find project root (go.mod)")
		}
		dir = parent
	}
}
