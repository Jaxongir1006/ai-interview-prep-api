package audit

import (
	"errors"
	"time"

	"go-enterprise-blueprint/internal/modules/audit/ctrl/asynctask"
	"go-enterprise-blueprint/internal/modules/audit/ctrl/cli"
	"go-enterprise-blueprint/internal/modules/audit/ctrl/consumer"
	"go-enterprise-blueprint/internal/modules/audit/ctrl/http"
	"go-enterprise-blueprint/internal/modules/audit/domain"
	"go-enterprise-blueprint/internal/modules/audit/embassy"
	"go-enterprise-blueprint/internal/modules/audit/infra/postgres"
	"go-enterprise-blueprint/internal/modules/audit/usecase"
	"go-enterprise-blueprint/internal/modules/audit/usecase/actionlog/getactionlogs"
	"go-enterprise-blueprint/internal/modules/audit/usecase/statuschangelog/getstatuschangelogs"
	"go-enterprise-blueprint/internal/portal"
	portalaudit "go-enterprise-blueprint/internal/portal/audit"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Consumers          consumer.Config `yaml:"consumers"`
	WorkerPollInterval time.Duration   `yaml:"worker_poll_interval" default:"1s"`
}

type Module struct {
	asynctaskCTRL *asynctask.Controller
	consumerCTRL  *consumer.Controller
	cliCTRL       *cli.Controller
	httpCTRL      *http.Controller

	portal portalaudit.Portal
}

func (m *Module) name() string {
	return "audit"
}

func New(
	cfg Config,
	brokerConfig kafka.BrokerConfig,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	var (
		err error
		m   = &Module{}
	)

	// Init repositories
	domainContainer := domain.NewContainer(
		postgres.NewActionLogRepo(dbConn, dbConn),
		postgres.NewStatusChangeLogRepo(dbConn, dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	// Init use cases
	usecaseContainer := usecase.NewContainer(
		getactionlogs.New(domainContainer),
		getstatuschangelogs.New(domainContainer),
	)

	// Init portal
	m.portal = embassy.New(domainContainer)

	// Init controllers
	m.cliCTRL = cli.NewController(usecaseContainer)
	m.httpCTRL = http.NewController(usecaseContainer, portalContainer, portalContainer.Auth(), httpServer)
	m.asynctaskCTRL, err = asynctask.NewController(dbConn, m.name(), cfg.WorkerPollInterval, usecaseContainer)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	m.consumerCTRL, err = consumer.NewController(cfg.Consumers, brokerConfig, usecaseContainer)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return m, nil
}

func (m *Module) Portal() portalaudit.Portal {
	return m.portal
}

func (m *Module) Start() error {
	var g errgroup.Group

	g.Go(m.asynctaskCTRL.Start)

	g.Go(m.consumerCTRL.Start)

	return errx.Wrap(g.Wait())
}

func (m *Module) Shutdown() error {
	errs := make(chan error, 2) // buffer size == controller count

	go func() { errs <- m.asynctaskCTRL.Shutdown() }()

	go func() { errs <- m.consumerCTRL.Shutdown() }()

	return errx.Wrap(errors.Join(<-errs, <-errs)) // <-errs count == controller count
}
