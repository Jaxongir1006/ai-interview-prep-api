# Code Style

**Style guides:** Effective Go, Go Code Review Comments, Uber Go Style Guide

## Naming Conventions

| Element      | Convention                                | Example                                 |
| ------------ | ----------------------------------------- | --------------------------------------- |
| Package      | lowercase, single-word, matches directory | `user`, `auth`, `adminlogin`            |
| File         | snake_case                                | `user_repository.go`, `create_order.go` |
| Operation ID | kebab-case                                | `admin-login`, `create-superadmin`      |

## Struct Design

- Exported fields first, then unexported
- Group related fields together
- Don't repeat default tags (e.g., omit `bun` tag if it defaults to lowercase)

## Comments

- **Do:** Start exported item comments with the item name
- **Do:** Keep comments up-to-date with code
- **Don't:** Write obvious comments (`// GetUser gets a user`)
- **Don't:** Write action/narrative comments (`// I fixed this`, `// Added this to handle...`, `// This was missing`, `// Changed this to work`)

## Error Handling Style

Always separate error-producing calls from error checks. Don't inline assignments in `if` statements.

```go
// Good
err := doSomething()
if err != nil {
    return errx.Wrap(err)
}

// Bad
if err := doSomething(); err != nil {
    return errx.Wrap(err)
}
```

## Formatting & Linting

We use **golangci-lint** for formatting and linting. Configuration is at `.golangci.yml` in the project root.
