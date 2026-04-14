package uow

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/answer"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/question"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/review"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/uowbase"
)

type Factory = uowbase.Factory[UnitOfWork]

type UnitOfWork interface {
	uowbase.UnitOfWork

	Session() session.Repo
	Question() question.Repo
	Answer() answer.Repo
	Review() review.Repo
}
