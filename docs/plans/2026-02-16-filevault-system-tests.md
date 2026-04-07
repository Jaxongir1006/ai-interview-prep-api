# Filevault System Tests Implementation Plan

> **For Codex:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement system tests for the filevault module's upload and download use cases.

**Architecture:** Tests interact with the running app through HTTP only — upload via multipart POST, download via GET. MinIO is used only by the background app, never by tests directly. The attach step (needed between upload and download) is a direct DB update via a state helper. No MinIO cleanup needed since upload paths are UUID-based and never collide.

**Tech Stack:** Go testing, httpexpect (multipart), bun (DB state helpers), testify/assert

---

### Task 1: Add `filevault` to database cleanup schemas

**Files:**
- Modify: `tests/state/database/empty.go:14-16`

**Step 1: Add filevault schema**

In `appSchemas()`, add `"filevault"` to the return slice:

```go
func appSchemas() []string {
	return []string{"auth", "audit", "taskmill", "alert", "filevault"}
}
```

**Step 2: Verify it compiles**

Run: `go build ./tests/...`
Expected: no errors

**Step 3: Commit**

```
feat: add filevault schema to test database cleanup
```

---

### Task 2: Create filevault state helpers — validate keys

**Files:**
- Create: `tests/state/filevault/validate.go`

**Step 1: Create the validate file**

```go
package filevault

//nolint:gochecknoglobals // static validation maps for test state
var (
	// validFileKeys defines the allowed keys for file test data.
	// Keys correspond to file.File entity fields.
	validFileKeys = map[string]bool{
		"id":               true,
		"original_name":    true,
		"stored_name":      true,
		"content_type":     true,
		"size":             true,
		"checksum":         true,
		"path":             true,
		"entity_type":      true,
		"entity_id":        true,
		"association_type":  true,
		"sort_order":       true,
		"uploaded_by":      true,
		"storage_status":   true,
	}
)
```

**Step 2: Verify it compiles**

Run: `go build ./tests/...`
Expected: no errors

**Step 3: Commit**

```
feat: add filevault test state validation keys
```

---

### Task 3: Create filevault state helpers — given.go

**Files:**
- Create: `tests/state/filevault/given.go`

**Step 1: Create GivenFiles helper**

This inserts file records directly into the DB. Defaults produce a "stored" file that's unattached. Tests for download errors (not attached, not ready) will override specific fields.

```go
package filevault

import (
	"fmt"
	"testing"
	"time"

	"go-enterprise-blueprint/internal/modules/filevault/domain/file"
	"go-enterprise-blueprint/internal/modules/filevault/infra/postgres"
	"go-enterprise-blueprint/pkg/anymap"
	"go-enterprise-blueprint/tests/state/database"

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

		f := &file.File{
			ID:              id,
			OriginalName:    anymap.String(d, "original_name", "test-file.png"),
			StoredName:      storedName,
			ContentType:     anymap.String(d, "content_type", "image/png"),
			Size:            cast.ToInt64(d["size"]),
			Checksum:        anymap.StringPtr(d, "checksum", lo.ToPtr("etag-test")),
			Path:            anymap.String(d, "path", fmt.Sprintf("2026/01/01/%s", storedName)),
			EntityType:      anymap.StringPtr(d, "entity_type", nil),
			EntityID:        anymap.Int64Ptr(d, "entity_id", nil),
			AssociationType: anymap.StringPtr(d, "association_type", nil),
			SortOrder:       cast.ToInt(d["sort_order"]),
			UploadedBy:      anymap.String(d, "uploaded_by", "test-actor"),
			StorageStatus:   anymap.String(d, "storage_status", file.StorageStatusStored),
		}

		if _, hasSize := d["size"]; !hasSize {
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
// This is a DB-only operation — it does not go through the portal/embassy.
func AttachFile(t *testing.T, fileID string, entityType string, entityID int64) {
	t.Helper()

	db := database.GetTestDB(t)

	ctx, cancel := database.QueryContext()
	defer cancel()

	_, err := db.NewUpdate().
		Model((*file.File)(nil)).
		ModelTableExpr("filevault.files").
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
```

**NOTE on `anymap.Int64Ptr`:** This function may not exist yet. Check if `anymap` has it. If not, use `cast` directly:

```go
// If anymap.Int64Ptr doesn't exist, use this pattern:
var entityID *int64
if v, ok := d["entity_id"]; ok && v != nil {
    id := cast.ToInt64(v)
    entityID = &id
}
```

**Step 2: Verify it compiles**

Run: `go build ./tests/...`
Expected: no errors (fix `anymap.Int64Ptr` if needed)

**Step 3: Commit**

```
feat: add filevault GivenFiles and AttachFile state helpers
```

---

### Task 4: Create filevault state helpers — get.go

**Files:**
- Create: `tests/state/filevault/get.go`

**Step 1: Create getter helpers**

```go
package filevault

import (
	"testing"

	"go-enterprise-blueprint/internal/modules/filevault/domain/file"
	"go-enterprise-blueprint/internal/modules/filevault/infra/postgres"
	"go-enterprise-blueprint/tests/state/database"
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
```

**Step 2: Verify it compiles**

Run: `go build ./tests/...`
Expected: no errors

**Step 3: Commit**

```
feat: add filevault getter state helpers
```

---

### Task 5: Write upload + download success system test

The main test: upload a real file via multipart POST, attach it via DB helper, then download it and verify the content matches.

**Files:**
- Create: `tests/system/modules/filevault/file/upload-download_test.go`

**Step 1: Write the test**

```go
//go:build system

package file_test

import (
	"net/http"
	"strings"
	"testing"

	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	statefilevault "go-enterprise-blueprint/tests/state/filevault"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/stretchr/testify/assert"
)

func TestUploadAndDownload_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t)

	fileContent := "hello filevault test content"
	fileName := "test-document.png"

	// WHEN — Upload
	uploadResp := trigger.UserAction(t).POST("/api/v1/files/upload").
		WithHeader("Authorization", "Bearer "+token).
		WithMultipart().
		WithFileBytes("file", fileName, []byte(fileContent)).
		Expect()

	// THEN — Upload response
	uploadResp.Status(http.StatusCreated)
	obj := uploadResp.JSON().Object()
	fileID := obj.Value("id").String().Raw()
	obj.Value("original_name").String().IsEqual(fileName)
	obj.Value("content_type").String().NotEmpty()
	obj.Value("size").Number().IsEqual(len(fileContent))

	assert.NotEmpty(t, fileID)

	// Verify DB state after upload
	f := statefilevault.GetFileByID(t, fileID)
	assert.Equal(t, "stored", f.StorageStatus)
	assert.NotNil(t, f.Checksum)
	assert.Nil(t, f.EntityID, "file should not be attached yet")

	// GIVEN — Attach file to an entity
	statefilevault.AttachFile(t, fileID, "test_entity", 1)

	// WHEN — Download
	downloadResp := trigger.UserAction(t).GET("/api/v1/files/download").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("id", fileID).
		Expect()

	// THEN — Download response
	downloadResp.Status(http.StatusOK)
	downloadResp.Header("Content-Disposition").Contains(fileName)
	downloadResp.Header("ETag").NotEmpty()

	body := downloadResp.Body().Raw()
	assert.Equal(t, fileContent, body)
}
```

**Step 2: Run the test**

Run: `make test-system` (or targeted: `go test -tags system -run TestUploadAndDownload_Success ./tests/system/modules/filevault/file/`)
Expected: PASS

**Step 3: Commit**

```
feat: add filevault upload + download success system test
```

---

### Task 6: Write upload error system tests

**Files:**
- Modify: `tests/system/modules/filevault/file/upload-download_test.go`

**Step 1: Write upload validation tests**

```go
func TestUpload_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		setupReq   func(req *httpexpect.Request) *httpexpect.Request
		wantStatus int
	}{
		{
			name: "missing file field",
			setupReq: func(req *httpexpect.Request) *httpexpect.Request {
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "file too large",
			setupReq: func(req *httpexpect.Request) *httpexpect.Request {
				// 11MB exceeds the 10MB limit in test.yaml
				bigContent := make([]byte, 11*1024*1024)
				return req.WithMultipart().
					WithFileBytes("file", "big.png", bigContent)
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := auth.GivenAuthToken(t)

			// WHEN
			req := trigger.UserAction(t).POST("/api/v1/files/upload").
				WithHeader("Authorization", "Bearer "+token)
			req = tc.setupReq(req)
			resp := req.Expect()

			// THEN
			resp.Status(tc.wantStatus)
		})
	}
}
```

**Step 2: Run the tests**

Run: `make test-system`
Expected: PASS

**Step 3: Commit**

```
feat: add filevault upload error system tests
```

---

### Task 7: Write download error system tests

**Files:**
- Modify: `tests/system/modules/filevault/file/upload-download_test.go`

**Step 1: Write download error tests**

```go
func TestDownload_Errors(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		// GIVEN
		database.Empty(t)
		token := auth.GivenAuthToken(t)

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/files/download").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("id", "00000000-0000-0000-0000-000000000000").
			Expect()

		// THEN
		resp.Status(http.StatusNotFound)
		resp.JSON().Object().Value("error").Object().
			Value("code").String().IsEqual("FILE_NOT_FOUND")
	})

	t.Run("file not attached", func(t *testing.T) {
		// GIVEN — file exists but is not attached to any entity
		database.Empty(t)
		token := auth.GivenAuthToken(t)
		files := statefilevault.GivenFiles(t, map[string]any{
			"storage_status": "stored",
		})

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/files/download").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("id", files[0].ID).
			Expect()

		// THEN
		resp.Status(http.StatusNotFound)
		resp.JSON().Object().Value("error").Object().
			Value("code").String().IsEqual("FILE_NOT_ATTACHED")
	})

	t.Run("file not ready", func(t *testing.T) {
		// GIVEN — file is attached but storage_status is pending
		database.Empty(t)
		token := auth.GivenAuthToken(t)
		files := statefilevault.GivenFiles(t, map[string]any{
			"storage_status": "pending",
			"entity_type":    "test_entity",
			"entity_id":      int64(1),
		})

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/files/download").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("id", files[0].ID).
			Expect()

		// THEN
		resp.Status(http.StatusBadRequest)
		resp.JSON().Object().Value("error").Object().
			Value("code").String().IsEqual("FILE_NOT_READY")
	})

	t.Run("missing file id", func(t *testing.T) {
		// GIVEN
		database.Empty(t)
		token := auth.GivenAuthToken(t)

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/files/download").
			WithHeader("Authorization", "Bearer "+token).
			Expect()

		// THEN
		resp.Status(http.StatusBadRequest)
	})
}
```

**Step 2: Run all system tests**

Run: `make test-system`
Expected: all pass

**Step 3: Commit**

```
feat: add filevault download error system tests
```

---

### Task 8: Write ETag / If-None-Match caching test

**Files:**
- Modify: `tests/system/modules/filevault/file/upload-download_test.go`

**Step 1: Write the caching test**

This extends the success flow — upload, attach, download to get ETag, then request again with `If-None-Match` and expect 304.

```go
func TestDownload_ETagCaching(t *testing.T) {
	// GIVEN — Upload and attach a file
	database.Empty(t)
	token := auth.GivenAuthToken(t)

	uploadResp := trigger.UserAction(t).POST("/api/v1/files/upload").
		WithHeader("Authorization", "Bearer "+token).
		WithMultipart().
		WithFileBytes("file", "cache-test.png", []byte("etag test content")).
		Expect()

	uploadResp.Status(http.StatusCreated)
	fileID := uploadResp.JSON().Object().Value("id").String().Raw()

	statefilevault.AttachFile(t, fileID, "test_entity", 1)

	// First download to get ETag
	firstResp := trigger.UserAction(t).GET("/api/v1/files/download").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("id", fileID).
		Expect()

	firstResp.Status(http.StatusOK)
	etag := firstResp.Header("ETag").Raw()
	assert.NotEmpty(t, etag)

	// WHEN — Request with matching If-None-Match
	cachedResp := trigger.UserAction(t).GET("/api/v1/files/download").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("If-None-Match", etag).
		WithQuery("id", fileID).
		Expect()

	// THEN — 304 Not Modified
	cachedResp.Status(http.StatusNotModified)
}
```

**Step 2: Run all system tests**

Run: `make test-system`
Expected: all pass

**Step 3: Commit**

```
feat: add filevault ETag caching system test
```

---

## Notes

- **Auth on filevault routes:** The upload UC calls `auth.MustUserContext(ctx)` to get the actor ID. Verify the routes have auth middleware applied. If not, upload/download calls without a token may panic instead of returning 401. If that's the case, auth middleware needs to be added to filevault routes before these tests will work properly.
- **`anymap.Int64Ptr`:** Check if this helper exists. If not, either add it to the `anymap` package or use an inline cast pattern in `GivenFiles`.
- **`entity_id` in GivenFiles:** The `entity_id` field is `*int64`. Use `cast.ToInt64` + pointer when the key is present in the data map.
