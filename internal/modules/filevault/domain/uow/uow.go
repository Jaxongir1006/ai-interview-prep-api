package uow

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain/file"
	"github.com/jaxongir1006/hire-ready-api/pkg/uowbase"
)

type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
	uowbase.UnitOfWork

	// Repository accessors
	File() file.Repo
}
