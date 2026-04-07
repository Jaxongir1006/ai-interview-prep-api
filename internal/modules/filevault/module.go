package filevault

import (
	"go-enterprise-blueprint/internal/modules/filevault/ctrl/http"
	"go-enterprise-blueprint/internal/modules/filevault/domain"
	"go-enterprise-blueprint/internal/modules/filevault/embassy"
	"go-enterprise-blueprint/internal/modules/filevault/infra/postgres"
	"go-enterprise-blueprint/internal/modules/filevault/usecase"
	"go-enterprise-blueprint/internal/modules/filevault/usecase/download"
	"go-enterprise-blueprint/internal/modules/filevault/usecase/upload"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/filevault"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/filestore/miniowr"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct {
	MaxFileSizeMB int64 `yaml:"max_file_size_mb" default:"10"`

	MinIO miniowr.Config `yaml:"minio" validate:"required"`
}

type Module struct {
	httpCTRL *http.Controller
	portal   filevault.Portal
}

func New(
	config Config,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	m := &Module{}

	fileStore, err := miniowr.New(config.MinIO)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	domainContainer := domain.NewContainer(
		fileStore,
		postgres.NewFileRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	usecaseContainer := usecase.NewContainer(
		upload.New(
			domainContainer,
			config.MaxFileSizeMB,
		),
		download.New(
			domainContainer,
		),
	)

	m.portal = embassy.New(
		domainContainer,
	)
	m.httpCTRL = http.NewController(usecaseContainer, portalContainer, httpServer)

	return m, nil
}

func (m *Module) Portal() filevault.Portal {
	return m.portal
}
