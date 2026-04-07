package embassy

import (
	"context"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain/file"

	"github.com/code19m/errx"
)

type detachRequest struct {
	EntityType string
	EntityID   int64
	AssocType  string
}

func (e *embassy) detach(ctx context.Context, req *detachRequest) error {
	// NO defer uow.DiscardUnapplied() — borrowed UOW must not rollback
	uow, err := e.domainContainer.UOWFactory().NewBorrowed(ctx)
	if err != nil {
		return errx.Wrap(err)
	}
	return errx.Wrap(uow.File().ClearEntityFields(ctx, file.Filter{
		EntityType:      &req.EntityType,
		EntityID:        &req.EntityID,
		AssociationType: &req.AssocType,
	}))
}
