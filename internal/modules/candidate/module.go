package candidate

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/ctrl/http"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/usecase/profile/completeonboarding"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct{}

type Module struct {
	httpCTRL *http.Controller
	portal   candidateportal.Portal
}

func New(
	_ Config,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	m := &Module{}

	domainContainer := domain.NewContainer(
		postgres.NewProfileRepo(dbConn),
		postgres.NewTopicPreferenceRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	usecaseContainer := usecase.NewContainer(
		completeonboarding.New(domainContainer),
	)

	m.portal = embassy.New(domainContainer)
	m.httpCTRL = http.NewController(usecaseContainer, portalContainer.Auth(), httpServer)

	return m, nil
}

func (m *Module) Portal() candidateportal.Portal {
	return m.portal
}
