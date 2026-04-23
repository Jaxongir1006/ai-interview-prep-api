package domain

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/answer"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/catalog"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/question"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/review"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain/uow"
)

type Container struct {
	sessionRepo  session.Repo
	questionRepo question.Repo
	answerRepo   answer.Repo
	reviewRepo   review.Repo
	catalogRepo  catalog.Repo
	uowFactory   uow.Factory
}

func NewContainer(
	sessionRepo session.Repo,
	questionRepo question.Repo,
	answerRepo answer.Repo,
	reviewRepo review.Repo,
	catalogRepo catalog.Repo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		sessionRepo:  sessionRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		reviewRepo:   reviewRepo,
		catalogRepo:  catalogRepo,
		uowFactory:   uowFactory,
	}
}

func (c *Container) SessionRepo() session.Repo {
	return c.sessionRepo
}

func (c *Container) QuestionRepo() question.Repo {
	return c.questionRepo
}

func (c *Container) AnswerRepo() answer.Repo {
	return c.answerRepo
}

func (c *Container) ReviewRepo() review.Repo {
	return c.reviewRepo
}

func (c *Container) CatalogRepo() catalog.Repo {
	return c.catalogRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
