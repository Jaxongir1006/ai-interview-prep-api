package auth

import (
	"errors"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/ctrl/asynctask"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/ctrl/cli"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/ctrl/consumer"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/ctrl/http"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/embassy"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/infra/postgres"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/createrole"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/deleterole"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getrolepermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getroles"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getuserpermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getuserroles"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/setrolepermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/setuserpermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/setuserroles"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/updaterole"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/cleanexpiredsessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/deletemysession"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/deletesession"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/deleteusersessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/getmysessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/getusersessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/adminlogin"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/changemypassword"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/createsuperadmin"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/createuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/disableuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/enableuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/getauthstats"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/getusers"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/logout"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/refreshtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/updateuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Consumers consumer.Config `yaml:"consumers"`

	AccessTokenTTL    time.Duration `yaml:"access_token_ttl"    default:"15m"`
	RefreshTokenTTL   time.Duration `yaml:"refresh_token_ttl"   default:"720h"` // 30 days
	MaxActiveSessions int           `yaml:"max_active_sessions" default:"5"`

	WorkerPollInterval time.Duration `yaml:"worker_poll_interval" default:"1s"`

	HashingCost int `yaml:"hashing_cost" default:"10" validate:"oneof=4 10 31"` // 4 for testing
}

type Module struct {
	asynctaskCTRL *asynctask.Controller
	consumerCTRL  *consumer.Controller
	cliCTRL       *cli.Controller
	httpCTRL      *http.Controller

	portal auth.Portal
}

func (m *Module) name() string {
	return "auth"
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
		postgres.NewUserRepo(dbConn),
		postgres.NewSessionRepo(dbConn),
		postgres.NewRoleRepo(dbConn),
		postgres.NewRolePermissionRepo(dbConn),
		postgres.NewUserRoleRepo(dbConn),
		postgres.NewUserPermissionRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	// Init use cases
	usecaseContainer := usecase.NewContainer(
		// Self-service
		adminlogin.New(
			domainContainer,
			portalContainer,
			cfg.AccessTokenTTL,
			cfg.RefreshTokenTTL,
			cfg.MaxActiveSessions,
		),
		refreshtoken.New(domainContainer, cfg.AccessTokenTTL, cfg.RefreshTokenTTL),
		logout.New(domainContainer, portalContainer),
		changemypassword.New(domainContainer, portalContainer, cfg.HashingCost),

		// User management
		createsuperadmin.New(domainContainer, cfg.HashingCost),
		getusers.New(domainContainer),
		getauthstats.New(domainContainer),
		createuser.New(domainContainer, portalContainer, cfg.HashingCost),
		updateuser.New(domainContainer, portalContainer, cfg.HashingCost),
		disableuser.New(domainContainer, portalContainer),
		enableuser.New(domainContainer, portalContainer),

		// Role management
		createrole.New(domainContainer, portalContainer),
		updaterole.New(domainContainer, portalContainer),
		deleterole.New(domainContainer, portalContainer),
		getroles.New(domainContainer),
		setrolepermissions.New(domainContainer, portalContainer),
		getrolepermissions.New(domainContainer),

		// Access management
		setuserroles.New(domainContainer, portalContainer),
		getuserroles.New(domainContainer),
		setuserpermissions.New(domainContainer, portalContainer),
		getuserpermissions.New(domainContainer),

		// Session management
		getmysessions.New(domainContainer),
		deletemysession.New(domainContainer, portalContainer),
		cleanexpiredsessions.New(domainContainer),
		getusersessions.New(domainContainer),
		deletesession.New(domainContainer, portalContainer),
		deleteusersessions.New(domainContainer, portalContainer),
	)

	// Init portal
	m.portal = embassy.New(domainContainer)

	// Init controllers
	m.cliCTRL = cli.NewController(usecaseContainer)
	m.httpCTRL = http.NewContoller(usecaseContainer, portalContainer, m.portal, httpServer)
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

func (m *Module) Portal() auth.Portal {
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

// --- CLI commands of auth module ---

func (m *Module) CreateSuperadmin(username, password string) error {
	var flags *cli.CreateSuperadminFlags
	if username != "" && password != "" {
		flags = &cli.CreateSuperadminFlags{
			Username: username,
			Password: password,
		}
	}
	return errx.Wrap(m.cliCTRL.CreateSuperadminCmd(flags))
}
