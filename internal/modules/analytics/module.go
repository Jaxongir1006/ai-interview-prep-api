package analytics

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/ctrl/http"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/pblc/dashboard"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getoverview"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getperformancetrend"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getrecentactivity"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getrecommendations"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/getstats"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/usecase/dashboard/gettopics"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	analyticsportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/analytics"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct{}

type Module struct {
	httpCTRL *http.Controller
	portal   analyticsportal.Portal
}

func New(
	_ Config,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	domainContainer := domain.NewContainer(
		postgres.NewProgressSummaryRepo(dbConn),
		postgres.NewTopicStatRepo(dbConn),
		postgres.NewAchievementDefinitionRepo(dbConn),
		postgres.NewCandidateAchievementRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	dashboardBuilder := dashboard.NewBuilder(domainContainer, portalContainer)
	usecaseContainer := usecase.NewContainer(
		getoverview.New(dashboardBuilder),
		getstats.New(dashboardBuilder),
		getperformancetrend.New(dashboardBuilder),
		gettopics.New(dashboardBuilder),
		getrecentactivity.New(dashboardBuilder),
		getrecommendations.New(dashboardBuilder),
	)

	return &Module{
		httpCTRL: http.NewController(usecaseContainer, portalContainer.Auth(), httpServer),
		portal:   embassy.New(domainContainer),
	}, nil
}

func (m *Module) Portal() analyticsportal.Portal {
	return m.portal
}
