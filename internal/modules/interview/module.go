package interview

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	interviewportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct{}

type Module struct {
	domainContainer *domain.Container
	portal          interviewportal.Portal
}

func New(
	_ Config,
	dbConn *bun.DB,
	_ *portal.Container,
	_ *server.HTTPServer,
) (*Module, error) {
	domainContainer := domain.NewContainer(
		postgres.NewSessionRepo(dbConn),
		postgres.NewQuestionRepo(dbConn),
		postgres.NewAnswerRepo(dbConn),
		postgres.NewReviewRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	return &Module{
		domainContainer: domainContainer,
		portal:          embassy.New(domainContainer),
	}, nil
}

func (m *Module) Portal() interviewportal.Portal {
	return m.portal
}
