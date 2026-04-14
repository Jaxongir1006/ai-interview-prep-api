package portal

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/analytics"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/filevault"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/interview"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/platform"
)

// Container holds every modules portal interface.
// It acts as a dependency injection container for the portal layer.
type Container struct {
	analytics analytics.Portal
	audit     audit.Portal
	auth      auth.Portal
	candidate candidate.Portal
	filevault filevault.Portal
	interview interview.Portal
	platform  platform.Portal
}

func (c *Container) SetAnalyticsPortal(analytics analytics.Portal) {
	c.analytics = analytics
}

func (c *Container) SetAuthPortal(auth auth.Portal) {
	c.auth = auth
}

func (c *Container) SetCandidatePortal(candidate candidate.Portal) {
	c.candidate = candidate
}

func (c *Container) SetAuditPortal(audit audit.Portal) {
	c.audit = audit
}

func (c *Container) SetFilevaultPortal(fv filevault.Portal) {
	c.filevault = fv
}

func (c *Container) SetInterviewPortal(interview interview.Portal) {
	c.interview = interview
}

func (c *Container) Auth() auth.Portal {
	return c.auth
}

func (c *Container) Analytics() analytics.Portal {
	return c.analytics
}

func (c *Container) Candidate() candidate.Portal {
	return c.candidate
}

func (c *Container) Audit() audit.Portal {
	return c.audit
}

func (c *Container) Filevault() filevault.Portal {
	return c.filevault
}

func (c *Container) Interview() interview.Portal {
	return c.interview
}

func (c *Container) SetPlatformPortal(platform platform.Portal) {
	c.platform = platform
}

func (c *Container) Platform() platform.Portal {
	return c.platform
}
