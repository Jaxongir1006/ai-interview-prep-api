package embassy

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/filevault/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/filevault"
)

type embassy struct {
	domainContainer *domain.Container
}

func New(
	domainContainer *domain.Container,
) filevault.Portal {
	return &embassy{
		domainContainer: domainContainer,
	}
}
