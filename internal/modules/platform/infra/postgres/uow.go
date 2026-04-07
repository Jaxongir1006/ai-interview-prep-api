package postgres

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/platform/domain/uow"
	"github.com/jaxongir1006/hire-ready-api/pkg/uowbase"
	"github.com/jaxongir1006/hire-ready-api/pkg/uowbase/pguowbase"

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
