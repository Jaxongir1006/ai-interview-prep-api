package embassy

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain/file"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/filevault"

	"github.com/code19m/errx"
)

func (e *embassy) Attach(ctx context.Context, req *filevault.AttachRequest) error {
	if len(req.FileIDs) == 0 {
		return nil
	}

	// NO defer uow.DiscardUnapplied() — borrowed UOW must not rollback
	uow, err := e.domainContainer.UOWFactory().NewBorrowed(ctx)
	if err != nil {
		return errx.Wrap(err)
	}

	// Find all files with SELECT FOR UPDATE
	fs, err := uow.File().ListForUpdate(ctx, file.Filter{
		IDs: req.FileIDs,
	})
	if err != nil {
		return errx.Wrap(err)
	}
	if len(fs) != len(req.FileIDs) {
		found := make(map[string]struct{}, len(fs))
		for _, f := range fs {
			found[f.ID] = struct{}{}
		}
		var missing []string
		for _, id := range req.FileIDs {
			if _, ok := found[id]; !ok {
				missing = append(missing, id)
			}
		}
		return errx.New("some files not found",
			errx.WithCode(filevault.CodeFileNotFound),
			errx.WithDetails(map[string]any{"missing_ids": missing}),
		)
	}

	// Validate all files
	for _, f := range fs {
		err = e.validateForAttach(f, req)
		if err != nil {
			return err
		}
	}

	// Update files
	fileMap := make(map[string]*file.File, len(fs))
	for i := range fs {
		fileMap[fs[i].ID] = &fs[i]
	}

	for i, id := range req.FileIDs {
		f := fileMap[id]
		f.EntityID = &req.EntityID
		f.EntityType = &req.EntityType
		f.AssociationType = &req.AssocType
		f.SortOrder = i + 1
	}

	err = uow.File().BulkUpdate(ctx, fs)
	return errx.Wrap(err)
}

func (e *embassy) validateForAttach(f file.File, req *filevault.AttachRequest) error {
	switch {
	case f.StorageStatus != file.StorageStatusStored:
		return errx.New("file is not ready for attachment",
			errx.WithCode(filevault.CodeFileNotReady),
			errx.WithDetails(errx.D{
				"file_id": f.ID,
			}),
		)
	case f.EntityID != nil:
		return errx.New("file is already attached",
			errx.WithCode(filevault.CodeFileAlreadyAttached),
			errx.WithDetails(errx.D{
				"file_id": f.ID,
			}),
		)
	case f.DeletedAt != nil:
		return errx.New("file has been deleted",
			errx.WithCode(filevault.CodeFileNotFound),
			errx.WithDetails(errx.D{
				"file_id": f.ID,
			}),
		)
	case req.ContentGroup != nil && !filevault.IsAllowedContentType(*req.ContentGroup, f.ContentType):
		allowedContentTypes := filevault.AllowedContentTypes(*req.ContentGroup)
		return errx.New("file content type not allowed",
			errx.WithCode(filevault.CodeFileTypeNotAllowed),
			errx.WithDetails(errx.D{
				"file_id":       f.ID,
				"content_type":  f.ContentType,
				"allowed_types": allowedContentTypes,
			}),
		)
	case req.MaxFileSizeMB != nil && f.Size > *req.MaxFileSizeMB*1024*1024:
		return errx.New("file size exceeds limit",
			errx.WithCode(filevault.CodeFileTooLarge),
			errx.WithDetails(errx.D{
				"file_id": f.ID,
				"size":    f.Size,
				"limit":   *req.MaxFileSizeMB * 1024 * 1024,
			}),
		)
	}
	return nil
}
