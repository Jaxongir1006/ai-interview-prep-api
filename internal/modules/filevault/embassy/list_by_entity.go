package embassy

import (
	"context"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain/file"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/filevault"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/sorter"
)

func (e *embassy) ListByEntity(ctx context.Context, req *filevault.ListByEntityRequest) ([]filevault.FileInfo, error) {
	fs, err := e.domainContainer.FileRepo().List(ctx, file.Filter{
		EntityType:      &req.EntityType,
		EntityID:        &req.EntityID,
		AssociationType: req.AssocType,
		SortOpts: sorter.SortOpts{
			sorter.Opt{F: "sort_order", D: sorter.Asc},
		},
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return toFileInfoSlice(fs), nil
}

func toFileInfoSlice(fs []file.File) []filevault.FileInfo {
	fis := make([]filevault.FileInfo, len(fs))
	for i, f := range fs {
		fis[i] = filevault.FileInfo{
			ID:           f.ID,
			OriginalName: f.OriginalName,
			ContentType:  f.ContentType,
			Size:         f.Size,
			AssocType:    f.AssociationType,
			SortOrder:    f.SortOrder,
		}
	}
	return fis
}
