package file

import (
	"context"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID              *string
	EntityType      *string
	EntityID        *int64
	AssociationType *string
	UploadedBy      *string
	StorageStatus   *string
	IDs             []string

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[File, Filter]

	ListForUpdate(ctx context.Context, filter Filter) ([]File, error)
	UpdateStorageStatus(ctx context.Context, id string, status string, checksum *string) error
	SoftDeleteByEntity(ctx context.Context, entityType string, entityID int64) error
	ClearEntityFields(ctx context.Context, filter Filter) error
}
