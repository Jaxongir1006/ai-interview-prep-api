package postgres

import (
	"go-enterprise-blueprint/internal/modules/platform/domain/uow"
	"go-enterprise-blueprint/pkg/uowbase"
	"go-enterprise-blueprint/pkg/uowbase/pguowbase"

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
