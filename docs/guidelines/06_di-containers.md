# DI Containers

Each layer (except controllers) provides a container struct that acts as a dependency injection container. Containers hold instances and expose them via getter methods. Wiring happens in `module.go` during application bootstrap.

## Containers by Layer

| Layer    | Container           | Holds                                                 |
| -------- | ------------------- | ----------------------------------------------------- |
| Domain   | `domain.Container`  | Repository interfaces, UOW factory, `pkg/` interfaces |
| PBLC     | `pblc.Container`    | PBLC component instances                              |
| Use Case | `usecase.Container` | Use case instances                                    |

Controllers don't provide containers — they are entry points, not dependencies.

## Domain Container

The foundational container. Every other layer receives it to access repositories and domain dependencies.

```go
type Container struct {
    adminRepo   user.AdminRepo
    sessionRepo session.Repo
    uowFactory  uow.Factory
}
```

### Shared `pkg/` Interfaces

If an interface already exists in the shared `pkg/` layer (e.g., a payment client, an SMS sender), add it directly to the domain container. Don't redefine or duplicate what `pkg/` already provides.

## Use Case Container

Holds all use case instances for the module. Controllers receive it to wire routes, consumers, and tasks.

```go
type Container struct {
    createSuperadmin     createsuperadmin.UseCase
    adminLogin           adminlogin.UseCase
    cleanExpiredSessions cleanexpiredsessions.UseCase
}
```

## PBLC Container

Holds all PBLC component instances. Use cases receive it alongside the domain container.

## Container Pattern

All containers follow the same structure:

- Unexported fields
- `NewContainer(...)` constructor accepting all dependencies
- Getter methods for each dependency

```go
type Container struct {
    adminRepo   user.AdminRepo
    sessionRepo session.Repo
    uowFactory  uow.Factory
}

func NewContainer(
    adminRepo user.AdminRepo,
    sessionRepo session.Repo,
    uowFactory uow.Factory,
) *Container {
    return &Container{
        adminRepo,
        sessionRepo,
        uowFactory,
    }
}

func (c *Container) AdminRepo() user.AdminRepo {
    return c.adminRepo
}

func (c *Container) SessionRepo() session.Repo {
    return c.sessionRepo
}

func (c *Container) UOWFactory() uow.Factory {
    return c.uowFactory
}
```
