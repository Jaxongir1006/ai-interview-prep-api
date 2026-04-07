package filevault

import (
	"fmt"
	"testing"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain/file"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/infra/postgres"
	"github.com/jaxongir1006/hire-ready-api/pkg/anymap"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// GivenFiles creates file records in the database for test setup.
// Each map in data represents a file with the following valid keys:
//   - id: string (default: generated UUID)
//   - original_name: string (default: "test-file.png")
//   - stored_name: string (default: "{id}.png")
//   - content_type: string (default: "image/png")
//   - size: int64 (default: 1024)
//   - checksum: *string (default: pointer to "etag-test")
//   - path: string (default: "2026/01/01/{stored_name}")
//   - entity_type: *string (default: nil)
//   - entity_id: *int64 (default: nil)
//   - association_type: *string (default: nil)
//   - sort_order: int (default: 0)
//   - uploaded_by: string (default: "test-actor")
//   - storage_status: string (default: "stored")
//
// Returns the created file entities.
func GivenFiles(t *testing.T, data ...map[string]any) []file.File {
	t.Helper()

	if len(data) == 0 {
		data = []map[string]any{{}}
	}

	db := database.GetTestDB(t)
	repo := postgres.NewFileRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	files := make([]file.File, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenFiles", validFileKeys, d)

		id := anymap.String(d, "id", uuid.NewString())
		storedName := anymap.String(d, "stored_name", id+".png")
		path := anymap.String(d, "path", fmt.Sprintf("2026/01/01/%s", storedName))

		var entityType *string
		if v, ok := d["entity_type"]; ok && v != nil {
			entityType = lo.ToPtr(cast.ToString(v))
		}

		var entityID *int64
		if v, ok := d["entity_id"]; ok && v != nil {
			entityID = lo.ToPtr(cast.ToInt64(v))
		}

		var associationType *string
		if v, ok := d["association_type"]; ok && v != nil {
			associationType = lo.ToPtr(cast.ToString(v))
		}

		f := &file.File{
			ID:              id,
			OriginalName:    anymap.String(d, "original_name", "test-file.png"),
			StoredName:      storedName,
			ContentType:     anymap.String(d, "content_type", "image/png"),
			Checksum:        anymap.StringPtr(d, "checksum", lo.ToPtr("etag-test")),
			Path:            path,
			EntityType:      entityType,
			EntityID:        entityID,
			AssociationType: associationType,
			SortOrder:       cast.ToInt(d["sort_order"]),
			UploadedBy:      anymap.String(d, "uploaded_by", "test-actor"),
			StorageStatus:   anymap.String(d, "storage_status", "stored"),
		}

		if _, hasSize := d["size"]; hasSize {
			f.Size = cast.ToInt64(d["size"])
		} else {
			f.Size = 1024
		}

		created, err := repo.Create(ctx, f)
		if err != nil {
			t.Fatalf("GivenFiles[%d]: failed to create file: %v", i, err)
		}

		files = append(files, *created)
	}

	return files
}

// AttachFile directly updates a file record to attach it to an entity.
// This is a DB-only operation -- it does not go through the portal/embassy.
func AttachFile(t *testing.T, fileID string, entityType string, entityID int64) {
	t.Helper()

	db := database.GetTestDB(t)

	ctx, cancel := database.QueryContext()
	defer cancel()

	_, err := db.NewUpdate().
		Model((*file.File)(nil)).
		ModelTableExpr("filevault.files AS file").
		Set("entity_type = ?", entityType).
		Set("entity_id = ?", entityID).
		Set("association_type = ?", "default").
		Set("sort_order = ?", 1).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", fileID).
		Exec(ctx)
	if err != nil {
		t.Fatalf("AttachFile: failed to attach file %q: %v", fileID, err)
	}
}
