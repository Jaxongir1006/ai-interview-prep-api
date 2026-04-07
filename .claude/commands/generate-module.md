# Generate Module

Scaffold a new module with all required layers, controllers, documentation, and application wiring.

The module name is provided as: $ARGUMENTS

## Steps

1. **Validate input.** The module name must be provided. It should be lowercase, single-word (e.g., `billing`, `audit`, `platform`).

2. **Ask the user** for the module's purpose — a brief description of what this module does. Use this to fill `overview.md`.

3. **Create the module directory structure.** Use the exact layout below, replacing `{module}` with the module name:

```
internal/modules/{module}/
├── module.go
├── domain/
│   └── container.go
├── usecase/
│   └── container.go
├── pblc/
│   └── .gitkeep
├── infra/
│   └── postgres/
│       └── const.go
├── ctrl/
│   ├── http/
│   │   └── http.go
│   ├── cli/
│   │   └── cli.go
│   ├── consumer/
│   │   └── consumer.go
│   └── asynctask/
│       └── asynctask.go
└── embassy/
    └── embassy.go
```

4. **Create portal interface.** Create `internal/portal/{module}/{module}.go` with an empty Portal interface.

5. **Create documentation.** Create `docs/specs/modules/{module}/`:
   - `overview.md` — fill Purpose and Responsibilities from user's description
   - `ERD.md` — empty ERD template
   - `usecases/` directory with `.gitkeep`

6. **Create test directories.** Create with `.gitkeep`:
   - `tests/system/modules/{module}/`
   - `tests/state/{module}/`

7. **Wire up in the application layer.** Make the following changes to existing files:

   **`internal/app/app.go`:**
   - Add import for the new module package
   - Add `{module} *{module}.Module` field to the `app` struct
   - Add `{Module} {module}.Config \`yaml:"{module}"\``to the`Config` struct under module configs section

   **`internal/app/run.go`:**
   - In `initModules()`: add module initialization (`a.{module}, err = {module}.New(...)`) after existing modules
   - In `initModules()`: add portal wiring (`portalContainer.Set{Module}Portal(a.{module}.Portal())`) in the portal section
   - In `runHighLevelComponents()`: add `g.Go(a.{module}.Start)`

   **`internal/app/shutdown.go`:**
   - In `shutdownHighLevelComponents()`: add shutdown item for the new module

   **`internal/portal/container.go`:**
   - Add import for the new module's portal package
   - Add field, setter, and getter for the new module's portal

## File Templates

Use the Go module path `go-enterprise-blueprint` for imports.

### module.go

```go
package {module}

import (
	"errors"
	"go-enterprise-blueprint/internal/modules/{module}/ctrl/asynctask"
	"go-enterprise-blueprint/internal/modules/{module}/ctrl/cli"
	"go-enterprise-blueprint/internal/modules/{module}/ctrl/consumer"
	"go-enterprise-blueprint/internal/modules/{module}/ctrl/http"
	"go-enterprise-blueprint/internal/modules/{module}/domain"
	"go-enterprise-blueprint/internal/modules/{module}/embassy"
	"go-enterprise-blueprint/internal/modules/{module}/usecase"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/{module}"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Consumers          consumer.Config `yaml:"consumers"`
	WorkerPollInterval time.Duration   `yaml:"worker_poll_interval" default:"1s"`
}

type Module struct {
	asynctaskCTRL *asynctask.Controller
	consumerCTRL  *consumer.Controller
	cliCTRL       *cli.Controller
	httpCTRL      *http.Controller

	portal {module}.Portal
}

func (m *Module) name() string {
	return "{module}"
}

func New(
	cfg Config,
	brokerConfig kafka.BrokerConfig,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	var (
		err error
		m   = &Module{}
	)

	// Init repositories
	domainContainer := domain.NewContainer()

	// Init use cases
	usecaseContainer := usecase.NewContainer()

	// Init portal
	m.portal = embassy.New()

	// Init controllers
	m.cliCTRL = cli.NewController(usecaseContainer)
	m.httpCTRL = http.NewController(usecaseContainer, portalContainer, httpServer)
	m.asynctaskCTRL, err = asynctask.NewController(dbConn, m.name(), cfg.WorkerPollInterval, usecaseContainer)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	m.consumerCTRL, err = consumer.NewController(cfg.Consumers, brokerConfig, usecaseContainer)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return m, nil
}

func (m *Module) Portal() {module}.Portal {
	return m.portal
}

func (m *Module) Start() error {
	var g errgroup.Group

	g.Go(m.asynctaskCTRL.Start)

	g.Go(m.consumerCTRL.Start)

	return errx.Wrap(g.Wait())
}

func (m *Module) Shutdown() error {
	errs := make(chan error, 2) // buffer size == controller count

	go func() { errs <- m.asynctaskCTRL.Shutdown() }()

	go func() { errs <- m.consumerCTRL.Shutdown() }()

	return errx.Wrap(errors.Join(<-errs, <-errs)) // <-errs count == controller count
}
```

### domain/container.go

```go
package domain

type Container struct{}

func NewContainer() *Container {
	return &Container{}
}
```

### usecase/container.go

```go
package usecase

type Container struct{}

func NewContainer() *Container {
	return &Container{}
}
```

### embassy/embassy.go

```go
package embassy

import "go-enterprise-blueprint/internal/portal/{module}"

type embassy struct{}

func New() {module}.Portal {
	return &embassy{}
}
```

### internal/portal/{module}/{module}.go

```go
package {module}

type Portal interface{}
```

### ctrl/http/http.go

```go
package http

import (
	"go-enterprise-blueprint/internal/modules/{module}/usecase"
	"go-enterprise-blueprint/internal/portal"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
)

type Controller struct {
	usecaseContainer *usecase.Container
	portalContainer  *portal.Container
}

func NewController(
	usecaseContainer *usecase.Container,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) *Controller {
	ctrl := &Controller{
		usecaseContainer,
		portalContainer,
	}

	httpServer.RegisterRouter(ctrl.initRoutes)
	return ctrl
}

func (c *Controller) initRoutes(r fiber.Router) {
	_ = r.Group("/api/v1/{module}")

	// Add your routes here...
}
```

### ctrl/cli/cli.go

```go
package cli

import (
	"go-enterprise-blueprint/internal/modules/{module}/usecase"
)

type Controller struct {
	usecaseContainer *usecase.Container
}

func NewController(usecaseContainer *usecase.Container) *Controller {
	return &Controller{
		usecaseContainer,
	}
}
```

### ctrl/consumer/consumer.go

Follow the exact same pattern as the auth module's consumer controller at `internal/modules/auth/ctrl/consumer/consumer.go`, replacing the module name. Start with all consumers commented out.

### ctrl/asynctask/asynctask.go

Follow the exact same pattern as the auth module's asynctask controller at `internal/modules/auth/ctrl/asynctask/asynctask.go`, replacing the module name. Start with empty `registerTasks()` and `registerSchedules()`.

### infra/postgres/const.go

```go
package postgres

const schemaName = "{module}"
```

### docs/specs/modules/{module}/overview.md

Fill this based on the user's description of the module's purpose:

```markdown
# {Module} Module

## Purpose

{Fill from user's description}

## Responsibilities

{Extract responsibilities from user's description}

## Domain Main Entities

| Entity | Description |
| ------ | ----------- |

See ERD.md for entity relationships.
```

### docs/specs/modules/{module}/ERD.md

````markdown
# {Module} Module ERD

```mermaid
erDiagram
```
````
