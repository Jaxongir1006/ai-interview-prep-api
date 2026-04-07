package postgres

import (
	"go-enterprise-blueprint/internal/modules/filevault/domain/file"
	"go-enterprise-blueprint/internal/modules/filevault/domain/uow"
	"go-enterprise-blueprint/pkg/uowbase/pguowbase"

	"github.com/uptrace/bun"
)

func NewUOWFactory(db *bun.DB) uow.Factory {
	return pguowbase.NewGenericFactory(
		db,
		schemaName,
		func(base *pguowbase.Base) uow.UnitOfWork {
			return &pgUOW{Base: base}
		},
	)
}

type pgUOW struct {
	*pguowbase.Base
}

func (u *pgUOW) File() file.Repo {
	return NewFileRepo(u.IDB())
}
