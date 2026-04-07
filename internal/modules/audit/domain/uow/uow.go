package uow

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/audit/domain/actionlog"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/audit/domain/statuschangelog"
	"github.com/jaxongir1006/hire-ready-api/pkg/uowbase"
)

// Factory defines an interface for creating new instances of the UnitOfWork.
type Factory = uowbase.Factory[UnitOfWork]

// UnitOfWork represents a single unit of work, typically mapping to a database transaction.
// It provides access to various repositories and methods to finalize or discard changes.
type UnitOfWork interface {
	uowbase.UnitOfWork

	// Repository accessors
	ActionLog() actionlog.Repo
	StatusChangeLog() statuschangelog.Repo
}
