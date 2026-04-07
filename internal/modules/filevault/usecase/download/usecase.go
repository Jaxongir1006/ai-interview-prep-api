package download

import (
	"context"
	"io"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain/file"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/filestore"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	FileID string
}

type Response struct {
	Body         io.ReadCloser
	ContentType  string
	Size         int64
	OriginalName string
	Checksum     *string
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
) UseCase {
	return &usecase{
		domainContainer: domainContainer,
	}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "download" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// File metadata
	f, err := uc.domainContainer.FileRepo().Get(ctx, file.Filter{
		ID: &in.FileID,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, file.CodeFileNotFound)
	}

	if f.EntityID == nil {
		return nil, errx.New("file is not attached",
			errx.WithCode(file.CodeFileNotAttached),
			errx.WithType(errx.T_NotFound),
		)
	}

	// Check storage status
	if f.StorageStatus != file.StorageStatusStored {
		return nil, errx.New("file is not ready for download",
			errx.WithCode(file.CodeFileNotReady),
			errx.WithType(errx.T_Validation),
		)
	}

	// Get file from storage
	obj, err := uc.domainContainer.FileStore().Get(ctx, f.Path)
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, filestore.CodeFileNotFound)
	}

	return &Response{
		Body:         obj.Content,
		ContentType:  f.ContentType,
		Size:         f.Size,
		OriginalName: f.OriginalName,
		Checksum:     f.Checksum,
	}, nil
}
