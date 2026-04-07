package uow

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain/file"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"
)

type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
	uowbase.UnitOfWork

	// Repository accessors
	File() file.Repo
}
