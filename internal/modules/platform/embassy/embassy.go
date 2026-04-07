package embassy

import "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/platform"

type embassy struct{}

func New() platform.Portal {
	return &embassy{}
}
