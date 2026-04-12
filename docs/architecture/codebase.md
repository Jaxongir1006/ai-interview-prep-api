# Codebase Architecture

## File Structure

```
project-root/
├── cmd/                    # Application entry points (Cobra CLI)
├── config/                 # Configuration files (YAML per environment)
├── docs/                   # Documentation
├── internal/               # Private application code
│   ├── app/                # Application bootstrap and lifecycle
│   ├── modules/            # Business modules (see Module Structure below)
│   └── portal/             # Cross-module communication interfaces
├── pkg/                    # Shared packages (project-specific utilities)
├── migrations/             # Database migrations (goose)
├── scripts/                # Utility scripts
└── tests/                  # System and integration tests
```

## Module Structure

Expanded view of `internal/modules/{module}/` — each module is a self-contained business capability:

```
internal/modules/{module}/
├── module.go              # Module initialization and wiring
├── domain/                # Domain layer
│   ├── container.go       # Domain container (DI)
│   └── {domain}/          # Grouped entities and infrastructure interfaces
├── usecase/               # Use case layer
│   └── {domain}/
│       └── {usecase}/
│           └── usecase.go
├── pblc/                  # Packaged Business Logic Components
├── infra/                 # Implementations of domain interfaces
│   ├── postgres/          # Repository implementations
│   └── http/              # External HTTP client implementations
├── ctrl/                  # Controller layer
│   ├── http/              # HTTP handlers
│   ├── cli/               # CLI commands
│   ├── consumer/          # Event consumers
│   └── asynctask/         # Async task workers, schedulers
└── embassy/               # Portal implementation for this module
```

## Backend Layers

### 1. Domain Layer (`domain/`)

Business entities, value objects, and repository interfaces.

**Responsibilities:**

- Define entity structures
- Define repository interfaces (contracts)
- Pure business concepts, no implementation details

### 2. Use Case Layer (`usecase/`)

Application-specific business logic.

**Responsibilities:**

- Implement business operations
- Orchestrate domain objects and repositories
- One use case = one business operation
- Transaction boundaries

### 3. PBLC Layer (`pblc/`)

Packaged Business Logic Components — reusable business logic shared across use cases.

**Responsibilities:**

- Deduplicate logic across multiple use cases
- Strict input validation
- Called only from use cases, never from controllers

### 4. Infrastructure Layer (`infra/`)

External system implementations.

**Responsibilities:**

- Repository implementations (PostgreSQL, Redis, etc.)
- External HTTP clients
- Third-party integrations

### 5. Controller Layer (`ctrl/`)

Entry points that adapt external interfaces to use cases.

**Responsibilities:**

- HTTP handlers, CLI commands, event consumers, async tasks
- Request parsing and response formatting
- One-to-one mapping to use cases
- No business logic — delegate to use cases

## Cross-Module Communication

Modules communicate **only** through Portal interfaces. No direct imports between modules.

- **Portal** (`internal/portal/`) — defines the interfaces each module exposes
- **Embassy** (`internal/modules/{module}/embassy/`) — implements the portal interface inside the module

```
internal/portal/                          # Portal interfaces (contracts)
├── container.go                          # Aggregates all module portals
├── auth/                                 # Auth module portal interface
└── filevault/                            # Filevault module portal interface

internal/modules/{module}/embassy/        # Embassy (portal implementation)
└── embassy.go
```

**Key principles:**

- Modules never import from other modules directly
- Portal exposes only necessary functionality
- Embassy implements the portal from the module's side
- No transaction sharing between modules
- Each module owns its data
- Embassies are wired to the portal container during application startup (`internal/app/`)
