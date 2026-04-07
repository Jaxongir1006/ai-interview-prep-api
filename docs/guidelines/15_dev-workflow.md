# Development Workflow

The document-first principle applies at every level of engineering. Each level feeds into the next:

| Level    | Document first                      | Then implement             |
| -------- | ----------------------------------- | -------------------------- |
| System   | Design flows in `docs/specs/flows/` | Break into modules and UCs |
| Module   | Write `overview.md` + `ERD.md`      | Break into use cases       |
| Use case | Write UC doc from template          | Follow the 5-phase cycle   |

The use case is the atomic unit of development. Every feature, change, or fix follows the same cycle: **Analyze → Document → Implement → Test → Review & Verify**. This applies equally to developers and AI agents.

## The Cycle

```
┌──────────┐     ┌──────────┐     ┌─────────────┐     ┌──────┐     ┌──────────────────┐
│ Analyze  │────▶│ Document │────▶│ Implement   │────▶│ Test │────▶│ Review & Verify  │
└──────────┘     └──────────┘     └─────────────┘     └──────┘     └──────────────────┘
                      ▲                  │
                      └──────────────────┘
                    feedback loop: update
                    docs when discovering
                    missing edge cases
```

## 1. Analyze

Before writing anything — docs or code — understand the context:

1. Ensure you have read the module's `overview.md` and `ERD.md`
2. Ensure you have read related existing use case documentations
3. Identify the use case type (UserAction, EventSubscriber, AsyncTask, ManualCommand)
4. Identify which module owns this feature
5. Check if cross-module communication is needed (portals)
6. Check if reusable logic already exists in PBLC

**Output**: clear understanding of what to build, where it fits, and what already exists.

## 2. Document

Documentation comes before implementation. Always. Follow the documentation guideline (`14_documentation.md`) strictly.

### New Use Case

1. Pick the correct template from `docs/specs/templates/usecases/`
2. Fill in all sections following `14_documentation.md`
3. Update module `overview.md` if adding new core entities or responsibilities
4. Update `ERD.md` and provide new migration file if schema changes are needed
5. Review from a higher level — are any related use cases missing?

### Modifying Existing Use Case

1. Read current use case documentation fully
2. Update the documentation **first** — inputs, outputs, execute steps
3. Then proceed to implementation

### Bug Fix

1. Read the use case documentation for the affected feature
2. Determine: is the documentation wrong, or is the implementation wrong?
3. Fix whichever is incorrect — documentation stays the source of truth

## 3. Implement

Build layer by layer, bottom to top. Follow the relevant guideline for each layer:

| Order | Layer          |
| ----- | -------------- |
| 1     | Migrations     |
| 2     | Domain         |
| 3     | Infrastructure |
| 4     | PBLC           |
| 5     | Use Case       |
| 6     | Controller     |
| 7     | DI Containers  |

### Rules

- Don't add logic that isn't in the documentation — update docs first, then code
- Don't skip layers — even if a layer seems trivial, keep the structure
- When you discover a missing edge case during implementation, go back to step 2 (Document) before coding it

## 4. Test

Write system tests after implementation. Follow the testing guideline (`11_testing.md`) for patterns, structure, and rules.

### Deriving Tests from Documentation

Use case documentation directly maps to test cases:

| Documentation element     | Test case                                                                               |
| ------------------------- | --------------------------------------------------------------------------------------- |
| Success execute steps     | One comprehensive success test                                                          |
| Simple failures           | One table-driven test — input varies, setup is the same (validation, missing fields)    |
| Complex failures          | Dedicated test per scenario — setup itself is different (inactive user, limit exceeded) |
| "Start UOW" / "Apply UOW" | Verify atomicity — partial failure = no state change                                    |

## 5. Review & Verify

Two distinct checklists applied together: **Verify** checks correctness and sync, **Review** checks code quality.

### Verify (sync checks)

- [ ] **Docs ↔ Code**: execute steps in docs match `Execute` method comments
- [ ] **Docs ↔ Tests**: every documented behavior has a corresponding test
- [ ] **Docs ↔ Docs**: usecase index in `README.md` is up-to-date
- [ ] **Guidelines**: implementation follows all relevant guidelines
- [ ] **Lint**: `make lint` passes
- [ ] **Tests**: `make test` and `make test-system` passes

### Review (quality checks)

- [ ] **Layer violations**: no business logic in controllers or repos, no infra concerns in use cases
- [ ] **Security**: no SQL injection in custom repo methods, no sensitive data (but not strict, e.g., `password`) in logs/error messages, auth checks in place
- [ ] **Simplification**: no unnecessary abstractions, no PBLC extraction for single-use logic, no over-engineered patterns
- [ ] **Performance**: no N+1 queries, UOW not held open during read-only operations, indexes exist for filtered/sorted columns
- [ ] **Portal discipline**: no direct cross-module imports, all cross-module communication goes through portals

### Scaling by Change Type

Not every change needs the full review checklist:

| Change Type     | Verify     | Review                                        |
| --------------- | ---------- | --------------------------------------------- |
| New use case    | Full       | Full                                          |
| New module      | Full       | Full                                          |
| Modify use case | Sync check | Layer violations, security, simplification    |
| Bug fix         | Sync check | Security (ensure fix doesn't introduce holes) |
| Refactor        | Sync check | Simplification, layer violations, performance |

## Workflow by Change Type

| Change Type     | Analyze           | Document                    | Implement                | Test              | Review & Verify       |
| --------------- | ----------------- | --------------------------- | ------------------------ | ----------------- | --------------------- |
| New use case    | Full analysis     | Write from template         | All layers               | New test file     | Full                  |
| Modify use case | Read existing     | Update docs first           | Changed layers           | Update tests      | Sync + partial review |
| Bug fix         | Read UC docs      | Fix docs if wrong           | Fix code                 | Add regression    | Sync + security       |
| New module      | Full analysis     | overview + ERD + UCs        | All layers + `module.go` | New test dir      | Full                  |
| Refactor        | Read all affected | Update if structure changes | Refactor code            | Ensure tests pass | Sync + partial review |

## The Feedback Loop

The documentation-implementation relationship is **bidirectional**:

```
Documentation ──────▶ Implementation
     (drives)

Implementation ──────▶ Documentation
   (informs updates)
```

- **Forward**: documentation drives what you implement
- **Backward**: if implementation reveals something docs missed (edge case, new error scenario, missing validation), update docs **before** coding the fix
- **Never**: implement undocumented behavior and leave it undocumented
