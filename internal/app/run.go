package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Jaxongir1006/ai-interview-prep-api/i18n"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/baseserver"

	"github.com/code19m/errx"
	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"github.com/rise-and-shine/pkg/pg"
	"github.com/rise-and-shine/pkg/rediswr"
	"golang.org/x/sync/errgroup"
)

func Run() error {
	app := newApp()
	defer app.shutdown()

	err := app.init()
	if err != nil {
		return errx.Wrap(err)
	}

	errChan := make(chan error)

	// run all high level components
	go func() {
		errChan <- app.runHighLevelComponents()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	// error occurred at module.Start
	case err = <-errChan:
		return errx.Wrap(err)

	// signal received, just return nil to trigger app.shutdown()
	case <-quit:
		return nil
	}
}

func (a *app) runHighLevelComponents() error {
	var g errgroup.Group

	g.Go(a.httpServer.Start)
	logger.With("address", a.cfg.HTTPServer.Address()).Info("HTTP server is running . . .")

	// Run your modules here...
	g.Go(a.auth.Start)
	g.Go(a.audit.Start)
	g.Go(a.platform.Start)

	return errx.Wrap(g.Wait())
}

func (a *app) init() error {
	err := a.initSharedComponents()
	if err != nil {
		return errx.Wrap(err)
	}

	err = a.migrateUp()
	if err != nil {
		return errx.Wrap(err)
	}

	err = a.initModules()
	return errx.Wrap(err)
}

func (a *app) initSharedComponents() error {
	var (
		err error
	)

	// set global meta infomations
	meta.SetServiceInfo(a.cfg.Service.Name, a.cfg.Service.Version)
	meta.SetLanguageMap(i18n.Translations, i18n.DefaultLang)

	// init logger
	logger.SetGlobal(a.cfg.Logger)

	// init db connection pool
	a.dbConn, err = pg.NewBunDB(a.cfg.Postgres)
	if err != nil {
		return errx.Wrap(err)
	}

	a.redisClient = rediswr.New(a.cfg.Redis)

	// init metrics
	// Metrics provider not implemented yet...

	// init tracing
	a.tracerShutdownFunc, err = tracing.InitGlobalTracer(a.cfg.Tracing)
	if err != nil {
		return errx.Wrap(err)
	}

	// init alerting
	err = alert.SetGlobal(a.cfg.Alert, a.dbConn)
	if err != nil {
		return errx.Wrap(err)
	}

	// init http server
	a.httpServer = baseserver.New(a.cfg.HTTPServer, a.cfg.CORS)
	a.httpServer.GetApp().Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })

	return nil
}

func (a *app) initModules() error {
	var (
		err error
	)

	portalContainer := &portal.Container{}

	// Init all your modules here...
	a.auth, err = auth.New(
		a.cfg.Auth, a.cfg.KafkaBroker, a.dbConn, a.redisClient, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}
	portalContainer.SetAuthPortal(a.auth.Portal())

	a.audit, err = audit.New(
		a.cfg.Audit, a.cfg.KafkaBroker, a.dbConn, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	a.analytics, err = analytics.New(
		a.cfg.Analytics, a.dbConn, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	a.candidate, err = candidate.New(
		a.cfg.Candidate, a.dbConn, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	// Filevault Module
	a.filevault, err = filevault.New(
		a.cfg.Filevault, a.dbConn, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	// Interview Module
	a.interview, err = interview.New(
		a.cfg.Interview, a.dbConn, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	// Platform
	a.platform, err = platform.New(
		a.cfg.Platform, a.cfg.KafkaBroker, a.dbConn, portalContainer, a.httpServer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	// Set all portal implementations here...
	portalContainer.SetAnalyticsPortal(a.analytics.Portal())
	portalContainer.SetAuditPortal(a.audit.Portal())
	portalContainer.SetCandidatePortal(a.candidate.Portal())
	portalContainer.SetFilevaultPortal(a.filevault.Portal())
	portalContainer.SetInterviewPortal(a.interview.Portal())
	portalContainer.SetPlatformPortal(a.platform.Portal())

	return nil
}
