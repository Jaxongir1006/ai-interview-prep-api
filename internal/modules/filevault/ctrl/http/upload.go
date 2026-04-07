package http

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/usecase/upload"

	"github.com/code19m/errx"
	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/val"
)

func (c *Controller) Upload(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return errx.New(
			"file is required",
			errx.WithType(errx.T_Validation),
			errx.WithCode(val.CodeValidationFailed),
		)
	}

	f, err := fileHeader.Open()
	if err != nil {
		return errx.New(
			"failed to read file",
			errx.WithType(errx.T_Validation),
			errx.WithCode(val.CodeValidationFailed),
		)
	}
	defer f.Close()

	resp, err := c.usecaseContainer.Upload().Execute(ctx.UserContext(), &upload.Request{
		File:         f,
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	return ctx.JSON(resp)
}
