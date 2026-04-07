package filevault

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/ctrl/http"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/embassy"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/infra/postgres"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/usecase"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/usecase/download"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/usecase/upload"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/filevault"

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
