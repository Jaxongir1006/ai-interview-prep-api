package uow

import "github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"

// Factory defines an interface for creating new instances of the UnitOfWork.
type Factory = uowbase.Factory[uowbase.UnitOfWork]
