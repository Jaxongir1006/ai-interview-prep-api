package http

import (
	"fmt"
	"io"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/usecase/download"

	"github.com/code19m/errx"
	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/val"
)

func (c *Controller) Download(ctx *fiber.Ctx) error {
	fileID := ctx.Query("id")
	if fileID == "" {
		return errx.New("id is required",
			errx.WithType(errx.T_Validation),
			errx.WithCode(val.CodeValidationFailed),
		)
	}

	resp, err := c.usecaseContainer.Download().Execute(ctx.UserContext(), &download.Request{
		FileID: fileID,
	})
	if err != nil {
		return errx.Wrap(err)
	}
	defer resp.Body.Close()

	// If-None-Match
	if resp.Checksum != nil && ctx.Get("If-None-Match") == *resp.Checksum {
		return errx.Wrap(
			ctx.SendStatus(fiber.StatusNotModified),
		)
	}

	// Set headers
	ctx.Set("Content-Type", resp.ContentType)
	ctx.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, resp.OriginalName))

	if resp.Checksum != nil {
		ctx.Set("ETag", *resp.Checksum)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errx.Wrap(err)
	}

	return errx.Wrap(ctx.Send(body))
}
