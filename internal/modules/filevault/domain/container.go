package domain

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain/file"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain/uow"

	"github.com/rise-and-shine/pkg/filestore"
)

type Container struct {
	fileStore  filestore.FileStore
	fileRepo   file.Repo
	uowFactory uow.Factory
}

func NewContainer(
	fileStore filestore.FileStore,
	fileRepo file.Repo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		fileStore:  fileStore,
		fileRepo:   fileRepo,
		uowFactory: uowFactory,
	}
}

func (c *Container) FileRepo() file.Repo {
	return c.fileRepo
}

func (c *Container) FileStore() filestore.FileStore {
	return c.fileStore
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
