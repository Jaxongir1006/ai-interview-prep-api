package uow

import (
	"go-enterprise-blueprint/internal/modules/audit/domain/actionlog"
	"go-enterprise-blueprint/internal/modules/audit/domain/statuschangelog"
	"go-enterprise-blueprint/pkg/uowbase"
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
