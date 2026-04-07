package uow

import "go-enterprise-blueprint/pkg/uowbase"

// Factory defines an interface for creating new instances of the UnitOfWork.
type Factory = uowbase.Factory[uowbase.UnitOfWork]
