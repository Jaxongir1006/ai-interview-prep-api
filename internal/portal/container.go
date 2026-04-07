package portal

import (
	"github.com/jaxongir1006/hire-ready-api/internal/portal/audit"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/esign"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/filevault"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/platform"
)

// Container holds every modules portal interface.
// It acts as a dependency injection container for the portal layer.
type Container struct {
	audit     audit.Portal
	auth      auth.Portal
	esign     esign.Portal
	filevault filevault.Portal
	platform  platform.Portal
}

func (c *Container) SetAuthPortal(auth auth.Portal) {
	c.auth = auth
}

func (c *Container) SetEsignPortal(esign esign.Portal) {
	c.esign = esign
}

func (c *Container) SetAuditPortal(audit audit.Portal) {
	c.audit = audit
}

func (c *Container) SetFilevaultPortal(fv filevault.Portal) {
	c.filevault = fv
}

func (c *Container) Auth() auth.Portal {
	return c.auth
}

func (c *Container) Esign() esign.Portal {
	return c.esign
}

func (c *Container) Audit() audit.Portal {
	return c.audit
}

func (c *Container) Filevault() filevault.Portal {
	return c.filevault
}

func (c *Container) SetPlatformPortal(platform platform.Portal) {
	c.platform = platform
}

func (c *Container) Platform() platform.Portal {
	return c.platform
}
