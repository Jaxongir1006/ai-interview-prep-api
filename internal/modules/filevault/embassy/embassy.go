package embassy

import (
	"github.com/jaxongir1006/hire-ready-api/internal/modules/filevault/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/filevault"
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
