package postgres

import (
	"context"
	"errors"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain/file"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

type fileRepo struct {
	idb bun.IDB

	repogen.Repo[file.File, file.Filter]
}

func NewFileRepo(idb bun.IDB) file.Repo {
	return &fileRepo{
		idb: idb,
		Repo: repogen.NewPgRepoBuilder[file.File, file.Filter](idb).
			WithSchemaName(schemaName).
			WithNotFoundCode(file.CodeFileNotFound).
			WithFilterFunc(fileFilterFunc).
			Build(),
	}
}

func fileFilterFunc(q *bun.SelectQuery, f file.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.EntityType != nil {
		q = q.Where("entity_type = ?", *f.EntityType)
	}
	if f.EntityID != nil {
		q = q.Where("entity_id = ?", *f.EntityID)
	}
	if f.AssociationType != nil {
		q = q.Where("association_type = ?", *f.AssociationType)
	}
	if f.UploadedBy != nil {
		q = q.Where("uploaded_by = ?", *f.UploadedBy)
	}
	if f.StorageStatus != nil {
		q = q.Where("storage_status = ?", *f.StorageStatus)
	}
	if f.IDs != nil {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}
	if f.Offset != nil {
		q = q.Offset(*f.Offset)
	}
	for _, o := range f.SortOpts {
		q = q.Order(o.ToSQL())
	}
	return q
}

func (r *fileRepo) ListForUpdate(ctx context.Context, filter file.Filter) ([]file.File, error) {
	var files []file.File

	q := r.idb.NewSelect().
		Model(&files).
		ModelTableExpr(schemaName + ".files AS file").
		For("UPDATE")

	q = fileFilterFunc(q, filter)

	err := q.Scan(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return files, nil
}

func (r *fileRepo) UpdateStorageStatus(ctx context.Context, id string, status string, checksum *string) error {
	q := r.idb.NewUpdate().
		Model((*file.File)(nil)).
		ModelTableExpr(schemaName+".files AS file").
		Set("storage_status = ?", status).
		Set("updated_at = now()").
		Where("id = ?", id)

	if checksum != nil {
		q = q.Set("checksum = ?", checksum)
	}

	_, err := q.Exec(ctx)
	return errx.Wrap(err)
}

func (r *fileRepo) SoftDeleteByEntity(ctx context.Context, entityType string, entityID int64) error {
	q := r.idb.NewUpdate().
		Model((*file.File)(nil)).
		ModelTableExpr(schemaName+".files AS file").
		Set("updated_at = now()").
		Set("deleted_at = now()").
		Where("entity_type = ?", entityType).
		Where("entity_id = ?", entityID).
		Where("deleted_at IS NULL")

	_, err := q.Exec(ctx)
	return errx.Wrap(err)
}

func (r *fileRepo) ClearEntityFields(ctx context.Context, filter file.Filter) error {
	if filter.EntityType == nil || filter.EntityID == nil {
		return errors.New("ClearEntityFields requires EntityType and EntityID filter")
	}

	q := r.idb.NewUpdate().
		Model((*file.File)(nil)).
		ModelTableExpr(schemaName + ".files AS file").
		Set("entity_type = NULL").
		Set("entity_id = NULL").
		Set("association_type = NULL").
		Set("sort_order = 0").
		Set("updated_at = now()").
		Where("deleted_at IS NULL")

	if filter.EntityType != nil {
		q = q.Where("entity_type = ?", *filter.EntityType)
	}
	if filter.EntityID != nil {
		q = q.Where("entity_id = ?", *filter.EntityID)
	}
	if filter.AssociationType != nil {
		q = q.Where("association_type = ?", *filter.AssociationType)
	}

	_, err := q.Exec(ctx)
	return errx.Wrap(err)
}
