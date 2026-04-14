package postgres

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/answer"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/question"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/review"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/uow"
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

func (u *pgUOW) Session() session.Repo {
	return NewSessionRepo(u.IDB())
}

func (u *pgUOW) Question() question.Repo {
	return NewQuestionRepo(u.IDB())
}

func (u *pgUOW) Answer() answer.Repo {
	return NewAnswerRepo(u.IDB())
}

func (u *pgUOW) Review() review.Repo {
	return NewReviewRepo(u.IDB())
}
