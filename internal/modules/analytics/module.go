package analytics

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	analyticsportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/analytics"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct{}

type Module struct {
	portal analyticsportal.Portal
}

func New(
	_ Config,
	dbConn *bun.DB,
	_ *portal.Container,
	_ *server.HTTPServer,
) (*Module, error) {
	domainContainer := domain.NewContainer(
		postgres.NewProgressSummaryRepo(dbConn),
		postgres.NewTopicStatRepo(dbConn),
		postgres.NewAchievementDefinitionRepo(dbConn),
		postgres.NewCandidateAchievementRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	return &Module{
		portal: embassy.New(domainContainer),
	}, nil
}

func (m *Module) Portal() analyticsportal.Portal {
	return m.portal
}
