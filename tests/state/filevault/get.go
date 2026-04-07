package filevault

import (
	"testing"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain/file"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/infra/postgres"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"
)

// GetFileByID retrieves a file by ID.
// Fails the test if the file is not found.
func GetFileByID(t *testing.T, id string) *file.File {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewFileRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	f, err := repo.Get(ctx, file.Filter{ID: &id})
	if err != nil {
		t.Fatalf("GetFileByID: failed to get file %q: %v", id, err)
	}

	return f
}

// FileExists checks if a file with the given ID exists.
func FileExists(t *testing.T, id string) bool {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewFileRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	exists, err := repo.Exists(ctx, file.Filter{ID: &id})
	if err != nil {
		t.Fatalf("FileExists: failed to check file %q: %v", id, err)
	}

	return exists
}

// FileCount returns the total number of file records.
func FileCount(t *testing.T) int {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewFileRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := repo.Count(ctx, file.Filter{})
	if err != nil {
		t.Fatalf("FileCount: failed to count files: %v", err)
	}

	return count
}
