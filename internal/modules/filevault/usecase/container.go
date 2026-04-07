package usecase

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/usecase/download"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/usecase/upload"
)

type Container struct {
	fileUpload   upload.UseCase
	fileDownload download.UseCase
}

func NewContainer(
	fileUpload upload.UseCase,
	fileDownload download.UseCase,
) *Container {
	return &Container{
		fileUpload:   fileUpload,
		fileDownload: fileDownload,
	}
}

func (c *Container) Upload() upload.UseCase {
	return c.fileUpload
}

func (c *Container) Download() download.UseCase {
	return c.fileDownload
}
