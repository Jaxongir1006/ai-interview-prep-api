package embassy

import (
	"context"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/filevault"

	"github.com/code19m/errx"
)

func (e *embassy) Replace(ctx context.Context, req *filevault.ReplaceRequest) error {
	// Detach old files
	err := e.detach(ctx, &detachRequest{
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		AssocType:  req.AssocType,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	// Attach new files
	err = e.Attach(ctx, &filevault.AttachRequest{
		FileIDs:       req.NewFileIDs,
		EntityType:    req.EntityType,
		EntityID:      req.EntityID,
		AssocType:     req.AssocType,
		ContentGroup:  req.ContentGroup,
		MaxFileSizeMB: req.MaxFileSizeMB,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}
