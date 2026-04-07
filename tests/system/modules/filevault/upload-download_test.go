//go:build system

package file_test

import (
	"net/http"
	"testing"

	"go-enterprise-blueprint/tests/state/auth"
	"go-enterprise-blueprint/tests/state/database"
	statefilevault "go-enterprise-blueprint/tests/state/filevault"
	"go-enterprise-blueprint/tests/system/trigger"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

// Minimal 1x1 white PNG.
//
//nolint:gochecknoglobals // static test data
var testPNGContent = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, // PNG signature
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // 8-bit RGB
	0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41, // IDAT chunk
	0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00, // compressed
	0x00, 0x00, 0x02, 0x00, 0x01, 0xe2, 0x21, 0xbc, // data
	0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, // IEND chunk
	0x44, 0xae, 0x42, 0x60, 0x82,
}

func TestUploadAndDownload_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenAuthToken(t)

	fileName := "test-document.png"

	// WHEN -- Upload
	uploadResp := trigger.UserAction(t).POST("/api/v1/filevault/upload").
		WithHeader("Authorization", "Bearer "+token).
		WithMultipart().
		WithFileBytes("file", fileName, testPNGContent).
		Expect()

	// THEN -- Upload response
	uploadResp.Status(http.StatusOK)
	obj := uploadResp.JSON().Object()
	fileID := obj.Value("id").String().Raw()
	obj.Value("original_name").String().IsEqual(fileName)
	obj.Value("content_type").String().IsEqual("image/png")
	obj.Value("size").Number().IsEqual(len(testPNGContent))

	assert.NotEmpty(t, fileID)

	// Verify DB state after upload
	f := statefilevault.GetFileByID(t, fileID)
	assert.Equal(t, "stored", f.StorageStatus)
	assert.NotNil(t, f.Checksum)
	assert.Nil(t, f.EntityID, "file should not be attached yet")

	// GIVEN -- Attach file to an entity
	statefilevault.AttachFile(t, fileID, "test_entity", 1)

	// WHEN -- Download
	downloadResp := trigger.UserAction(t).GET("/api/v1/filevault/download").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("id", fileID).
		Expect()

	// THEN -- Download response
	downloadResp.Status(http.StatusOK)
	downloadResp.Header("Content-Disposition").Contains(fileName)
	downloadResp.Header("ETag").NotEmpty()

	body := downloadResp.Body().Raw()
	assert.Equal(t, string(testPNGContent), body)
}

func TestFilevault_Unauthorized(t *testing.T) {
	database.Empty(t)

	// Upload without token
	trigger.UserAction(t).POST("/api/v1/filevault/upload").
		WithMultipart().
		WithFileBytes("file", "test.png", testPNGContent).
		Expect().
		Status(http.StatusUnauthorized)

	// Download without token
	trigger.UserAction(t).GET("/api/v1/filevault/download").
		WithQuery("id", "00000000-0000-0000-0000-000000000000").
		Expect().
		Status(http.StatusUnauthorized)
}

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
				// Add PNG header so MIME type passes
				copy(bigContent, testPNGContent[:8])
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
			req := trigger.UserAction(t).POST("/api/v1/filevault/upload").
				WithHeader("Authorization", "Bearer "+token)
			req = tc.setupReq(req)
			resp := req.Expect()

			// THEN
			resp.Status(tc.wantStatus)
		})
	}
}

func TestDownload_Errors(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		// GIVEN
		database.Empty(t)
		token := auth.GivenAuthToken(t)

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/filevault/download").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("id", "00000000-0000-0000-0000-000000000000").
			Expect()

		// THEN
		resp.Status(http.StatusNotFound)
		resp.JSON().Object().Value("error").Object().
			Value("code").String().IsEqual("FILE_NOT_FOUND")
	})

	t.Run("file not attached", func(t *testing.T) {
		// GIVEN -- file exists but is not attached to any entity
		database.Empty(t)
		token := auth.GivenAuthToken(t)
		files := statefilevault.GivenFiles(t, map[string]any{
			"storage_status": "stored",
		})

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/filevault/download").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("id", files[0].ID).
			Expect()

		// THEN
		resp.Status(http.StatusNotFound)
		resp.JSON().Object().Value("error").Object().
			Value("code").String().IsEqual("FILE_NOT_ATTACHED")
	})

	t.Run("file not ready", func(t *testing.T) {
		// GIVEN -- file is attached but storage_status is pending
		database.Empty(t)
		token := auth.GivenAuthToken(t)
		files := statefilevault.GivenFiles(t, map[string]any{
			"storage_status": "pending",
			"entity_type":    "test_entity",
			"entity_id":      int64(1),
		})

		// WHEN
		resp := trigger.UserAction(t).GET("/api/v1/filevault/download").
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
		resp := trigger.UserAction(t).GET("/api/v1/filevault/download").
			WithHeader("Authorization", "Bearer "+token).
			Expect()

		// THEN
		resp.Status(http.StatusBadRequest)
	})
}

func TestDownload_ETagCaching(t *testing.T) {
	// GIVEN -- Upload and attach a file
	database.Empty(t)
	token := auth.GivenAuthToken(t)

	uploadResp := trigger.UserAction(t).POST("/api/v1/filevault/upload").
		WithHeader("Authorization", "Bearer "+token).
		WithMultipart().
		WithFileBytes("file", "cache-test.png", testPNGContent).
		Expect()

	uploadResp.Status(http.StatusOK)
	fileID := uploadResp.JSON().Object().Value("id").String().Raw()

	statefilevault.AttachFile(t, fileID, "test_entity", 1)

	// First download to get ETag
	firstResp := trigger.UserAction(t).GET("/api/v1/filevault/download").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("id", fileID).
		Expect()

	firstResp.Status(http.StatusOK)
	etag := firstResp.Header("ETag").Raw()
	assert.NotEmpty(t, etag)

	// WHEN -- Request with matching If-None-Match
	cachedResp := trigger.UserAction(t).GET("/api/v1/filevault/download").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("If-None-Match", etag).
		WithQuery("id", fileID).
		Expect()

	// THEN -- 304 Not Modified
	cachedResp.Status(http.StatusNotModified)
}
