package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase/pguowbase"

	"github.com/uptrace/bun"
)

func NewUOWFactory(db *bun.DB) uow.Factory {
	return pguowbase.NewGenericFactory(
		db,
		"platform",
		func(base *pguowbase.Base) uowbase.UnitOfWork {
			return base
		},
	)
}
