package embassy

import "github.com/jaxongir1006/hire-ready-api/internal/portal/platform"

type embassy struct{}

func New() platform.Portal {
	return &embassy{}
}
