package uow

import (
	"go-enterprise-blueprint/internal/modules/filevault/domain/file"
	"go-enterprise-blueprint/pkg/uowbase"
)

type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
	uowbase.UnitOfWork

	// Repository accessors
	File() file.Repo
}
