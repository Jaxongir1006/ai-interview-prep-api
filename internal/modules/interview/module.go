package interview

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/ctrl/http"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/usecase/catalog/getonboardingoptions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	interviewportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct{}

type Module struct {
	domainContainer *domain.Container
	httpCTRL        *http.Controller
	portal          interviewportal.Portal
}

func New(
	_ Config,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	domainContainer := domain.NewContainer(
		postgres.NewSessionRepo(dbConn),
		postgres.NewQuestionRepo(dbConn),
		postgres.NewAnswerRepo(dbConn),
		postgres.NewReviewRepo(dbConn),
		postgres.NewCatalogRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)
	portalImpl := embassy.New(domainContainer)

	usecaseContainer := usecase.NewContainer(
		getonboardingoptions.New(portalImpl),
	)

	return &Module{
		domainContainer: domainContainer,
		httpCTRL:        http.NewController(usecaseContainer, portalContainer.Auth(), httpServer),
		portal:          portalImpl,
	}, nil
}

func (m *Module) Portal() interviewportal.Portal {
	return m.portal
}
