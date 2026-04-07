package upload

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain/file"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/filevault"

	"github.com/code19m/errx"
	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	File         io.Reader
	OriginalName string
	Size         int64
}

type Response struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	maxFileSizeMB int64,
) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		maxFileSize:     maxFileSizeMB << 20,
	}
}

type usecase struct {
	domainContainer *domain.Container
	maxFileSize     int64
}

func (uc *usecase) OperationID() string { return "upload" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Resolve actor ID from context
	actorID := auth.MustUserContext(ctx).UserID

	// Validate size
	if in.Size > uc.maxFileSize {
		return nil, errx.New(
			fmt.Sprintf("max file size: %dMB", uc.maxFileSize>>20),
			errx.WithType(errx.T_Validation),
			errx.WithCode(file.CodeFileTooLarge),
		)
	}

	// Detect and validate mime type
	mimeType, reader, err := detectMimeType(in.File)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if !filevault.IsAllowedContentType(filevault.ContentGroupAll, mimeType) {
		return nil, errx.New(
			fmt.Sprintf("file content type not allowed: %s", mimeType),
			errx.WithType(errx.T_Validation),
			errx.WithCode(file.CodeFileTypeNotAllowed),
		)
	}

	// Generate file metadata
	fileID := uuid.NewString()
	ext := filepath.Ext(in.OriginalName)
	storedName := fileID + ext
	path := fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), storedName)

	// DB insert (pending)
	f, err := uc.domainContainer.FileRepo().Create(ctx, &file.File{
		ID:            fileID,
		OriginalName:  in.OriginalName,
		StoredName:    storedName,
		ContentType:   mimeType,
		Size:          in.Size,
		Path:          path,
		StorageStatus: file.StorageStatusPending,
		UploadedBy:    actorID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Upload to storage
	info, err := uc.domainContainer.FileStore().Upload(ctx, path, reader)
	if err != nil {
		// Best-effort status update; primary error takes precedence
		_ = uc.domainContainer.FileRepo().UpdateStorageStatus(ctx, f.ID, file.StorageStatusFailed, nil)
		return nil, errx.Wrap(err)
	}

	// Update status to stored
	err = uc.domainContainer.FileRepo().UpdateStorageStatus(ctx, f.ID, file.StorageStatusStored, &info.ETag)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		ContentType:  f.ContentType,
		Size:         f.Size,
	}, nil
}
