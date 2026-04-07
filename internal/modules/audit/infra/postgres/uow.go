package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/actionlog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/statuschangelog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase/pguowbase"

	"github.com/uptrace/bun"
)

func NewUOWFactory(db *bun.DB) uow.Factory {
	return pguowbase.NewGenericFactory(
		db,
		schemaName,
		func(base *pguowbase.Base) uow.UnitOfWork {
			return &pgUOW{Base: base, db: db}
		},
	)
}

type pgUOW struct {
	*pguowbase.Base
	db *bun.DB
}

func (u *pgUOW) ActionLog() actionlog.Repo {
	return NewActionLogRepo(u.IDB(), u.db)
}

func (u *pgUOW) StatusChangeLog() statuschangelog.Repo {
	return NewStatusChangeLogRepo(u.IDB(), u.db)
}
