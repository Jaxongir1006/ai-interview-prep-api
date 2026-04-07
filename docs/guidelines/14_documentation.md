# Documentation

This section outlines best practices for writing and maintaining documentation.

## Document-First Approach

- Documentation serves as the source of truth
- Don't code any use-case code before documenting
- Implementation follows documentation
- Keep documentation and implementation consistent
- Update docs **before** or **with** code changes
- Fix inconsistencies immediately when found

## Documentation Structure

```
docs/
├── architecture/          # System architecture (codebase structure, layers)
├── guidelines/            # Development guidelines (this folder)
├── misc/                  # Miscellaneous (research PDFs, notes)
└── specs/                 # Specifications
    ├── modules/           # Module documentation
    │   └── {module}/
    │       ├── overview.md
    │       ├── ERD.md
    │       └── usecases/
    │           └── {domain}/
    │               └── {operation}.md
    ├── flows/             # Business and technical flows
    │   └── {flow-name}.md
    └── templates/         # Documentation templates
        └── usecases/
```

## Module Documentation

Each module has its own folder in `docs/specs/modules/{module}/` with the following structure:

### overview.md (required)

Module purpose, responsibilities and domain main entities. Concise and clear.

### ERD.md (required)

Entity relationship diagram using Mermaid. Defines all database tables with column types and relationships.

- Follow database normalization rules — unnormalized ERD is bad design
- Denormalization is acceptable only when done **intentionally** for performance reasons, never by accident
- Describe fields where the name alone doesn't convey the business meaning — skip universally understood fields (`id`, `created_at`, `updated_at`)

```
admins {
    UUID id PK
    VARCHAR username UK
    VARCHAR password_hash
    BOOLEAN is_active "whether admin can log in"
    TIMESTAMPTZ last_active_at "updated on each login"
    TIMESTAMPTZ created_at
    TIMESTAMPTZ updated_at
}
```

### usecases/ (required)

Use case documentation organized by domains:

```
usecases/
├── user/
│   ├── admin-login.md
│   └── create-superadmin.md
└── session/
    └── clean-expired-sessions.md
```

## Use Case Documentation

Use case docs serve as **API specifications** — they replace Swagger/OpenAPI. A frontend developer should be able to implement their side purely from the use case document, without asking backend any questions. Every input, output, validation rule and behavior must be explicit.

Just follow the template strictly — fill in every section, nothing extra needed. Execution steps should describe **what** the system does, not **how** it does it.

Do NOT include: logging, metrics, query mechanics, infrastructure concerns.

### Transaction Boundaries in Execution Steps

When a use case requires atomic writes, explicitly mark the transaction boundary using **Start UOW** and **Apply UOW** steps:

```
- Validate input
- Find user by ID
- Start UOW                          ← transaction begins
- Update user record
- Create audit log entry
- Apply UOW                          ← transaction commits
- Return updated user
```

This makes atomicity visible at the documentation level — a reader can immediately see which operations are grouped into a single transaction. Steps between "Start UOW" and "Apply UOW" either all succeed or all roll back.

- Place **Start UOW** after all read-only checks and validations — don't hold a transaction open while doing lookups that don't need it
- Place **Apply UOW** after the last write operation, before returning results
- If a use case has no multi-write atomicity needs, omit both steps entirely — single writes don't need UOW

### Analysis Process

When defining a new feature or use case:

**1. Understand the problem.** Analyze existing use cases and flows. Ask clarifying questions if necessary:

- What is the business goal?
- Who are the actors (users, systems)?
- What triggers this operation?
- What is the expected outcome?
- What are the constraints and business rules?

**2. Classify the use case type:**

| Choose             | When                                                                                                |
| ------------------ | --------------------------------------------------------------------------------------------------- |
| `user_action`      | User initiates action and waits for response (API endpoints, CRUD, queries)                         |
| `event_subscriber` | Triggered by something that happened elsewhere, no immediate response (notifications, side effects) |
| `async_task`       | Runs at specific times/intervals or as background jobs (reports, cleanup, asynchronous tasks)       |
| `manual_command`   | Operator runs manually via CLI, one-time or maintenance operations                                  |

**3. Document the use case.** Use the appropriate template from `docs/specs/templates/usecases/`.

**4. Review from a higher level.** Ensure no necessary use cases are missing for the business flow being developed.

### Templates

| Use Case Type      | Template Path                                       |
| ------------------ | --------------------------------------------------- |
| `user_action`      | `docs/specs/templates/usecases/user_action.md`      |
| `async_task`       | `docs/specs/templates/usecases/async_task.md`       |
| `event_subscriber` | `docs/specs/templates/usecases/event_subscriber.md` |
| `manual_command`   | `docs/specs/templates/usecases/manual_command.md`   |

### Documenting Output Fields

Output sections must list **every field** the API actually returns. If the return type is a domain entity, include all serialized fields — including BaseModel timestamps (`created_at`, `updated_at`) and foreign keys — not just the "interesting" ones. A frontend developer relies on this as the complete API contract.

### Documenting Sort Fields

When a use case accepts dynamic sort (a `sort` query parameter), the documentation **must** list the allowed sortable fields and the default sort. A bare `sort: string, optional` is incomplete — a frontend developer cannot know which fields are valid without this information.

**Format:** `sort`: string, optional — sortable fields: {comma-separated list}. Default: {field}:{direction}

**Example:**

```
- `sort`: string, optional — sortable fields: created_at, code, service, operation. Default: created_at:desc
```

### Documenting Validation Rules

- Always specify data types
- Define min/max lengths for strings
- Specify allowed values for enums
- Document format requirements (email, UUID, etc.)
- Note required vs optional fields
- Always specify nullable response fields with comment // nullable

### Naming Conventions

- **Operation ID** — kebab-case: `create-superadmin`, `admin-login`
- **Doc file name** — kebab-case: `admin-login.md`, `create-superadmin.md`
- **Grouping** — by module domains: `usecases/user/`, `usecases/session/`
- **Use verb-noun format** — `create-user`, `update-role`, `send-notification`
- **Be specific** — `create-admin-user` not just `create-user`

### Checklist Before Finalizing

- Use case type is correctly identified
- All inputs are documented with validation rules
- Dynamic sort fields: allowed fields and default are listed (not bare `sort: string, optional`)
- All outputs are documented
- Authorization requirements are clear
- Error scenarios are comprehensive
- Side effects are listed
- Business rules are explicit
- File is placed in correct directory

## Use Case Index (README)

Every use case must be registered in the `README.md` Use Case Index table. This serves as the central map of all system operations.

| Column      | Description                                                                                                         |
| ----------- | ------------------------------------------------------------------------------------------------------------------- |
| Use Case    | `module:operation-id` in backticks                                                                                  |
| Type        | Use case type. Plain `user_action` for human actors. Append `:sa` only for service account actors: `user_action:sa` |
| Permissions | Required permission(s) in backticks, or `-` if none (unauthenticated, CLI commands, async tasks)                    |
| Docs        | `[spec](...)` link to the use case doc                                                                              |

## Flows

Flows document business or technical processes that span multiple use cases — e.g., authentication flow, approval flow, tender flow. They live in `docs/specs/flows/`.

Unlike use cases, flows have **no fixed template**. Choose the format that best describes the flow:

- **Sequence diagrams** (Mermaid) — for flows with clear interactions between actors and systems
- **Plain Markdown** — for simpler flows that can be described as a series of steps
- **Mixed** — combine diagrams with explanatory text as needed

## Writing Style

- Be concise and precise
- Use active voice
- Avoid ambiguity
- Use Mermaid for diagrams
- Use JSON for input/output examples
- Use bullet points for execution steps
