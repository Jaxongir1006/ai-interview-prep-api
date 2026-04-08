package embassy

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/analytics/domain"
	analyticsportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/analytics"
)

type embassy struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) analyticsportal.Portal {
	return &embassy{
		domainContainer: domainContainer,
	}
}
