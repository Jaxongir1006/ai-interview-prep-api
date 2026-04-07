package embassy

import (
	"go-enterprise-blueprint/internal/modules/filevault/domain"
	"go-enterprise-blueprint/internal/portal/filevault"
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
