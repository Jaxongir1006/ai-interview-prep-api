# Testing

## Unit Tests — `pkg/` Layer

Unit tests are required for the `pkg/` layer — HTTP clients, shared utilities, and any reusable components. These packages are self-contained and benefit from isolated testing.

- HTTP clients are tested using the `net/http/httptest` package.
- For unit testing always create test files and package name of test file with `_test.go` suffix.

## System Tests

System tests cover the rest of the application. All use cases must have 100% system test coverage.

### Directory Structure

```
tests/
├── state/                    # Test state management
│   ├── database/             # DB helpers (GetTestDB, Empty)
│   └── {module}/             # Module state helpers (Given*, Get*)
└── system/
    ├── trigger/              # UC trigger helpers
    └── modules/{module}/{domain}/{operation}_test.go
```

### GIVEN-WHEN-THEN Pattern

- **GIVEN** — setup: `database.Empty(t)`, `Given*` functions, prepare payloads
- **WHEN** — execute: exactly ONE action via trigger helpers
- **THEN** — verify: assert response, verify DB state, check side effects

```go
func TestUseCaseName(t *testing.T) {
    // GIVEN
    database.Empty(t)
    admins := auth.GivenAdmins(t, map[string]any{})

    // WHEN
    resp := trigger.UserAction(t).POST("/api/v1/endpoint").
        WithJSON(payload).Expect()

    // THEN
    resp.Status(http.StatusOK)
    assert.Equal(t, expected, actual)
}
```

### Trigger Functions

| Use Case Type      | Trigger                                      |
| ------------------ | -------------------------------------------- |
| `user_action`      | `trigger.UserAction(t).POST(...)`            |
| `manual_command`   | `trigger.ManualCommand(t, args...)`          |
| `async_task`       | `trigger.AsyncTask(t, queue, opID, payload)` |
| `event_subscriber` | `trigger.EventSubscriber(t, topic, event)`   |

### State Helpers

#### Given Functions

```go
// Create with defaults
auth.GivenAdmins(t, map[string]any{})

// Create with custom fields
auth.GivenAdmins(t, map[string]any{"username": "custom", "is_active": false})

// Create multiple
auth.GivenAdmins(t, map[string]any{"username": "alice"}, map[string]any{"username": "bob"})
```

#### Getter Functions

```go
auth.GetAdminByUsername(t, "alice")           // Get entity by field
auth.AdminExists(t, "alice")                  // Check existence
auth.SessionCount(t, adminID)                 // Count related entities
auth.HasPermission(t, "admin", id, "perm")    // Check permissions
```

#### Passwords

Use pre-computed hashes to avoid bcrypt overhead:

- `auth.TestPassword1` — valid password in requests

### Test Rules

- **Isolation** — call `database.Empty(t)` at start of each test
- **Independence** — tests must not depend on other tests' data
- **Deterministic** — tests must produce same results every run
- **Table-driven** — always prefer table-driven tests where possible
- **Comprehensive success case** — cover everything related to success in a single test case. Don't split success checks across multiple tests. One test should verify the full outcome: no errors, correct response, correct DB state, side effects (e.g., password hashed, permissions created).
- **Minimize test count** — aim to cover the most cases with the fewest tests. Each test should be meaningful and verify as much as possible.

### Running Tests

```bash
make test        # Run unit tests
make test-system # Run system tests
```
