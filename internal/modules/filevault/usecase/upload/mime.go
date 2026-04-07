package upload

import (
	"bytes"
	"io"
	"net/http"

	"github.com/code19m/errx"
)

func detectMimeType(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 512)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", nil, errx.Wrap(err)
	}

	mimeType := http.DetectContentType(buf[:n])
	reader := io.MultiReader(bytes.NewReader(buf[:n]), r)

	return mimeType, reader, nil
}
