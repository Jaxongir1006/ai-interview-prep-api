package uow

import "github.com/jaxongir1006/hire-ready-api/pkg/uowbase"

// Factory defines an interface for creating new instances of the UnitOfWork.
type Factory = uowbase.Factory[uowbase.UnitOfWork]
