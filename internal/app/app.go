package app

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/platform"
	"github.com/Jaxongir1006/ai-interview-prep-api/pkg/baseserver"

	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"github.com/rise-and-shine/pkg/pg"
	"github.com/uptrace/bun"
)

type Config struct {
	// --- Shared configs ---

	Service ServiceConfig `yaml:"service" validate:"required"`

	Logger logger.Config `yaml:"logger" validate:"required"`

	Tracing tracing.Config `yaml:"tracing" validate:"required"`

	Alert alert.Config `yaml:"alert" validate:"required"`

	Postgres pg.Config `yaml:"postgres" validate:"required"`

	KafkaBroker kafka.BrokerConfig `yaml:"kafka_broker" validate:"required"`

	HTTPServer server.Config `yaml:"http_server" validate:"required"`

	CORS baseserver.CORSConfig `yaml:"cors"`

	// --- Module specific configs ---

	Auth auth.Config `yaml:"auth"`

	Audit audit.Config `yaml:"audit"`

	Analytics analytics.Config `yaml:"analytics"`

	Candidate candidate.Config `yaml:"candidate"`

	Filevault filevault.Config `yaml:"filevault"`

	Interview interview.Config `yaml:"interview"`

	Platform platform.Config `yaml:"platform"`
}

type app struct {
	cfg Config

	dbConn             *bun.DB
	tracerShutdownFunc func() error

	httpServer *server.HTTPServer

	auth      *auth.Module
	audit     *audit.Module
	analytics *analytics.Module
	candidate *candidate.Module
	filevault *filevault.Module
	interview *interview.Module
	platform  *platform.Module
}

func newApp() *app {
	app := &app{
		cfg: cfgloader.MustLoad[Config](),
	}
	return app
}

type ServiceConfig struct {
	// Name is the name of the service
	Name string `json:"name" validate:"required"`

	// Version is the version of the service
	Version string `json:"version" validate:"required"`
}
