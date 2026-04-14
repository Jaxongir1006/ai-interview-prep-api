package embassy

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/interview/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"
)

type embassy struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) interview.Portal {
	return &embassy{domainContainer: domainContainer}
}
