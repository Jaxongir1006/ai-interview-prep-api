package embassy

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/candidate/domain"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"
)

type embassy struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) candidateportal.Portal {
	return &embassy{
		domainContainer: domainContainer,
	}
}
