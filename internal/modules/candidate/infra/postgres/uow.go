package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/profile"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/topicpreference"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase/pguowbase"

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

func (u *pgUOW) Profile() profile.Repo {
	return NewProfileRepo(u.IDB())
}

func (u *pgUOW) TopicPreference() topicpreference.Repo {
	return NewTopicPreferenceRepo(u.IDB())
}
