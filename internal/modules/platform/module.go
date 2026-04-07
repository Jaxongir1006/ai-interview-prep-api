package platform

import (
	"errors"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/ctrl/asynctask"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/ctrl/cli"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/ctrl/consumer"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/ctrl/http"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/alerterror/cleanuperrors"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/alerterror/geterror"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/alerterror/geterrorstats"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/alerterror/listerrors"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/cleanupresults"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/getqueuestats"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/listdlqtasks"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/listqueues"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/listschedules"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/listtaskresults"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/purgedlq"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/purgequeue"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/requeuefromdlq"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform/usecase/taskmill/triggerschedule"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	portalplatform "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/platform"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/rise-and-shine/pkg/taskmill/console"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Consumers          consumer.Config `yaml:"consumers"`
	WorkerPollInterval time.Duration   `yaml:"worker_poll_interval" default:"1s"`
	AlertSchema        string          `yaml:"alert_schema"         default:"alert"`
}

type Module struct {
	asynctaskCTRL *asynctask.Controller
	consumerCTRL  *consumer.Controller
	cliCTRL       *cli.Controller
	httpCTRL      *http.Controller

	portal portalplatform.Portal
}

func (m *Module) name() string {
	return "platform"
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

	// Init console
	tmConsole, err := console.New(dbConn)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Init alert error repo
	alertSchema := cfg.AlertSchema
	if alertSchema == "" {
		alertSchema = "alert"
	}
	alertErrorRepo := postgres.NewAlertErrorRepo(dbConn, alertSchema)

	// Init UOW factory
	uowFactory := postgres.NewUOWFactory(dbConn)

	// Init domain container
	domainContainer := domain.NewContainer(tmConsole, alertErrorRepo, uowFactory)

	// Init use cases
	usecaseContainer := usecase.NewContainer(
		listqueues.New(domainContainer),
		getqueuestats.New(domainContainer),
		listdlqtasks.New(domainContainer),
		listtaskresults.New(domainContainer),
		listschedules.New(domainContainer),
		requeuefromdlq.New(domainContainer, portalContainer),
		purgequeue.New(domainContainer, portalContainer),
		purgedlq.New(domainContainer, portalContainer),
		cleanupresults.New(domainContainer, portalContainer),
		triggerschedule.New(domainContainer, portalContainer),
		listerrors.New(domainContainer),
		geterror.New(domainContainer),
		geterrorstats.New(domainContainer),
		cleanuperrors.New(domainContainer, portalContainer),
	)

	// Init portal
	m.portal = embassy.New()

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

func (m *Module) Portal() portalplatform.Portal {
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
