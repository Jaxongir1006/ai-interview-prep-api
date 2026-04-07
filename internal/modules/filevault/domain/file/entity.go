package file

import (
	"context"
	"time"

	"github.com/rise-and-shine/pkg/pg"
	"github.com/uptrace/bun"
)

const (
	CodeFileNotFound       = "FILE_NOT_FOUND"
	CodeFileTooLarge       = "FILE_TOO_LARGE"
	CodeFileTypeNotAllowed = "FILE_TYPE_NOT_ALLOWED"
	CodeFileUploadFailed   = "FILE_UPLOAD_FAILED"
	CodeFileNotReady       = "FILE_NOT_READY"
	CodeFileForbidden      = "FILE_FORBIDDEN"
	CodeFileNotAttached    = "FILE_NOT_ATTACHED"
)

const (
	StorageStatusPending = "pending"
	StorageStatusStored  = "stored"
	StorageStatusFailed  = "failed"
)

type File struct {
	pg.BaseModel

	ID           string  `json:"id"            bun:"id,pk"`
	OriginalName string  `json:"original_name" bun:"original_name,notnull"`
	StoredName   string  `json:"-"             bun:"stored_name,notnull"`
	ContentType  string  `json:"content_type"  bun:"content_type,notnull"`
	Size         int64   `json:"size"          bun:"size,notnull"`
	Checksum     *string `json:"-"             bun:"checksum"`
	Path         string  `json:"-"             bun:"path,notnull"`

	EntityType      *string `json:"-" bun:"entity_type"`
	EntityID        *int64  `json:"-" bun:"entity_id"`
	AssociationType *string `json:"-" bun:"association_type"`
	SortOrder       int     `json:"-" bun:"sort_order,notnull,default:0"`

	UploadedBy    string `json:"-" bun:"uploaded_by,notnull"`
	StorageStatus string `json:"-" bun:"storage_status,notnull"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time  `json:"updated_at" bun:"updated_at,notnull,default:now()"`
	DeletedAt *time.Time `json:"-"          bun:"deleted_at,soft_delete"`
}

func (f *File) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		f.CreatedAt = time.Now()
		f.UpdatedAt = f.CreatedAt
	case *bun.UpdateQuery:
		f.UpdatedAt = time.Now()
	}
	return nil
}
