# HireReady API

Backend API for HireReady, an AI-assisted interview preparation platform.

The project is built as a modular Go backend with document-first development, strict module boundaries, session-backed authentication, RBAC, candidate onboarding, file storage, audit logging, analytics foundations, and platform operations APIs.

## What Exists Today

| Area | Status |
| --- | --- |
| Authentication | Admin login, public login, email verification, OAuth login, refresh tokens, logout |
| Authorization | RBAC with roles, role permissions, direct user permissions, and permission middleware |
| Candidate profile | Registration-created profile, onboarding state, selected interview target role, experience level, and preferred topics |
| Analytics | Database model for progress summaries, topic stats, and achievements |
| File storage | Upload/download through Filevault with MinIO metadata and attachment support |
| Audit | Action logs and status-change logs |
| Platform operations | Taskmill queue/admin APIs and alert error browsing/cleanup |

The Interview module now defines the persistence foundation for sessions, questions, answers, and review outcomes. Use cases for starting interviews, submitting answers, and running AI reviews still need to be documented before implementation.

## Architecture

The codebase follows a module-first architecture:

```text
internal/
├── app/          # Application bootstrap and lifecycle
├── modules/      # Business modules
└── portal/       # Cross-module contracts
```

Each module owns its data and exposes cross-module behavior only through a portal interface. Modules do not import each other directly.

Typical module layout:

```text
internal/modules/{module}/
├── domain/       # Entities, value objects, repository interfaces
├── usecase/      # One package per business operation
├── pblc/         # Reusable business logic components
├── infra/        # Repository/client implementations
├── ctrl/         # HTTP, CLI, consumers, async tasks
└── embassy/      # Portal implementation
```

Implementation order is bottom-up:

1. Migrations
2. Domain
3. Infra
4. PBLC
5. Use case
6. Controller
7. DI/container wiring
8. Tests and verification

## Modules

| Module | Responsibility |
| --- | --- |
| `auth` | Identity, sessions, OAuth accounts, email verification, RBAC |
| `candidate` | Interview-prep profile data, onboarding state, selected interview catalog keys |
| `analytics` | Derived candidate metrics, topic performance, achievements |
| `interview` | Interview option catalogs, sessions, questions shown, answers, review outcomes |
| `filevault` | Object storage, file metadata, downloads, entity attachments |
| `audit` | Centralized user action and status-change logs |
| `platform` | Operational APIs for queues, schedules, task results, and errors |

## API Style

This API is operation-based, not REST-based.

- `GET` is used for queries.
- `POST` is used for mutations.
- No path parameters. Use query parameters for `GET` and JSON bodies for `POST`.
- List responses are wrapped in `{ "content": [] }`.
- Paginated responses include `page_number`, `page_size`, `count`, and `content`.
- Every response includes `X-Trace-ID`.

Examples:

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/auth/get-me
POST /api/v1/me/complete-onboarding
```

See [docs/specs/api/general.md](docs/specs/api/general.md) for response contracts.

## Documentation

Documentation is the source of truth.

```text
docs/
├── architecture/        # Codebase architecture
├── guidelines/          # Engineering rules and patterns
├── specs/
│   ├── api/             # Shared API contracts
│   ├── flows/           # Business flows
│   ├── modules/         # Module specs, ERDs, use cases
│   └── templates/       # Use case templates
└── plans/               # Implementation plans
```

Before implementing a module feature, read:

1. `docs/specs/modules/{module}/overview.md`
2. `docs/specs/modules/{module}/ERD.md`
3. Relevant files in `docs/specs/flows/`
4. The specific use case document only when working on that use case

Use case documents are API specs. They must describe every input, output field, validation rule, business step, error, and transaction boundary.

## Local Development

Copy local configuration:

```bash
cp config/local.yaml.example config/local.yaml
```

Start infrastructure and run the app:

```bash
make run
```

The default local services are:

| Service | URL |
| --- | --- |
| API | `http://localhost:9876` |
| PostgreSQL | `localhost:5432` |
| Redis | `localhost:6379` |
| Mailpit UI | `http://localhost:8025` |
| MinIO API | `http://localhost:9000` |
| MinIO Console | `http://localhost:9001` |
| AKHQ | `http://localhost:8080` |

Run only infrastructure:

```bash
make infra-up
```

Stop infrastructure:

```bash
make infra-down
```

## Database Migrations

Run migrations:

```bash
make migrate-up
```

Rollback one migration:

```bash
make migrate-down
```

Create a migration:

```bash
make migrate-create
```

Migration files live in [migrations](migrations).

## Verification

Run formatting:

```bash
make fmt
```

Run lint:

```bash
make lint
```

Run unit tests for shared packages:

```bash
make test
```

Run system tests:

```bash
make test-system
```

Before delivering backend code, the expected verification set is:

```bash
make lint
make test
make test-system
```

If lint fails, run `make fmt` and then `make lint` again.

## AI Interview Engine Direction

The current docs prepare the foundation for an AI interview product. The Interview module owns the first interview option catalogs plus persisted sessions, questions, answers, and reviews.

Recommended next modules:

| Module | Owns |
| --- | --- |
| `questionbank` | Future reusable/manual question catalog and reusable question metadata |
| `interview` | Interview target-role, experience-level, and topic catalogs; sessions; selected questions; submitted answers; timing; completion state; raw review outcomes |
| `analytics` | Aggregated progress derived from completed interviews and reviews |

Questions should be persisted. Even if AI generates them, the system should store the exact question shown to the user, its topic, difficulty, prompt/version metadata, and whether it is reusable or session-specific.

AI reviews should also be persisted. Store the submitted answer, the AI evaluation result, score, rubric breakdown, feedback, model/provider metadata, and timestamps. This keeps progress dashboards reproducible and avoids changing historical results when prompts or models change.

## Development Rules

- Document first, then implement.
- Keep docs, code, and tests in sync.
- Put business logic in use cases or PBLCs, not controllers or repositories.
- Use portals for cross-module communication.
- Do not create cross-module imports between `internal/modules/*`.
- Use UOW only when multiple writes need atomicity.
- Add or update system tests for use case behavior.

See [docs/guidelines](docs/guidelines) for the detailed engineering rules.
