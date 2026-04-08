package candidate

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct{}

type Module struct {
	portal candidateportal.Portal
}

func New(
	_ Config,
	dbConn *bun.DB,
	_ *portal.Container,
	_ *server.HTTPServer,
) (*Module, error) {
	domainContainer := domain.NewContainer(
		postgres.NewProfileRepo(dbConn),
		postgres.NewTopicPreferenceRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	return &Module{
		portal: embassy.New(domainContainer),
	}, nil
}

func (m *Module) Portal() candidateportal.Portal {
	return m.portal
}
