package filevault

import (
	"context"
)

const (
	CodeFileNotFound        = "FILE_NOT_FOUND"
	CodeFileTooLarge        = "FILE_TOO_LARGE"
	CodeFileTypeNotAllowed  = "FILE_TYPE_NOT_ALLOWED"
	CodeFileAlreadyAttached = "FILE_ALREADY_ATTACHED"
	CodeFileNotReady        = "FILE_NOT_READY"
)

type FileInfo struct {
	ID           string
	OriginalName string
	ContentType  string
	Size         int64
	AssocType    *string
	SortOrder    int
}

type AttachRequest struct {
	FileIDs       []string
	EntityType    string
	EntityID      int64
	AssocType     string
	ContentGroup  *ContentGroup // which content types are allowed
	MaxFileSizeMB *int64
}

type ListByEntityRequest struct {
	EntityType string
	EntityID   int64
	AssocType  *string
}

type ReplaceRequest struct {
	NewFileIDs    []string
	EntityType    string
	EntityID      int64
	AssocType     string
	ContentGroup  *ContentGroup
	MaxFileSizeMB *int64
}

type Portal interface {
	Attach(ctx context.Context, req *AttachRequest) error
	Replace(ctx context.Context, req *ReplaceRequest) error
	ListByEntity(ctx context.Context, req *ListByEntityRequest) ([]FileInfo, error)
	DeleteByEntity(ctx context.Context, entityType string, entityID int64) error
}
