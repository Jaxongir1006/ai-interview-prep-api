package embassy

import "go-enterprise-blueprint/internal/portal/platform"

type embassy struct{}

func New() platform.Portal {
	return &embassy{}
}
