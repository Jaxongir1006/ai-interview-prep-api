package embassy

import (
	"context"

	"github.com/code19m/errx"
)

func (e *embassy) DeleteByEntity(ctx context.Context, entityType string, entityID int64) error {
	return errx.Wrap(e.domainContainer.FileRepo().SoftDeleteByEntity(ctx, entityType, entityID))
}
