# Error Handling

Use `github.com/code19m/errx` package for all error handling.

## Core Rules

- Always handle errors, never ignore them
- Use `errx.Wrap()` to preserve stack traces
- Use error codes for programmatic error handling
- Never use `panic`

## Layer Responsibilities

### Types

Only the **use case layer** assigns error types. It knows the caller (actor) and decides what type of error to return (validation, not found, etc.). All downstream layers — PBLC, repository, infra — return errors with the default `errx.T_Internal` type. They never set error types themselves.

### Codes

Downstream layers **may** provide error codes when callers need to distinguish between different error scenarios. But don't overuse codes — only add a code when the caller genuinely needs to branch on it. If nobody needs to differentiate that error, just `errx.Wrap(err)` and move on.

### Summary

| Layer      | Error Types                         | Error Codes                                   |
| ---------- | ----------------------------------- | --------------------------------------------- |
| Controller | Handled by framework automatically  | —                                             |
| Use Case   | Assigns types (knows the actor)     | Checks codes from downstream, defines its own |
| PBLC       | Never sets types (default internal) | May return codes when callers need them       |
| Infra      | Never sets types (default internal) | May return codes (e.g., not found, conflict)  |

The controller layer maps error types to HTTP status codes, gRPC codes, etc. This logic is already implemented in the framework — developers don't need to handle it manually.

## Error Codes

Define as package-level constants in the domain layer:

```go
const (
    CodeAdminNotFound  = "ADMIN_NOT_FOUND"
    CodeIncorrectCreds = "INCORRECT_CREDENTIALS"
)
```

### Checking Codes

Use `errx.IsCodeIn` to branch on error codes from downstream layers:

```go
if errx.IsCodeIn(err, user.CodeAdminNotFound) {
    // do some logic
}
```

### Reassigning Types by Code — `errx.WrapWithTypeOnCodes`

When a downstream layer returns an error with a domain code, the use case often needs to change its type. `errx.WrapWithTypeOnCodes` does this in one call:

```go
return errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeAdminNotFound)
```

This is typically needed when a repository call is made **with user-provided input** — for example, looking up a user by an ID the caller supplied. If that ID doesn't exist, it's a not-found error for the caller.

If the same repository call is made with **internal data** (e.g., an ID fetched from another query), a not-found result would be an internal error (a data inconsistency). In that case, just `errx.Wrap(err)` — don't reassign the type.

### Inline Return with Wrap Functions

All `errx` wrap functions (`Wrap`, `WrapWithTypeOnCodes`) return `nil` when the input error is `nil`. When wrapping is the **last operation** before returning, skip the `if err != nil` check and return directly:

```go
// Good — clean and concise
role, err := repo.Create(ctx, &rbac.Role{Name: in.Name})
return role, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, rbac.CodeRoleNameConflict)

// Bad — unnecessary nil check
role, err := repo.Create(ctx, &rbac.Role{Name: in.Name})
if err != nil {
    return nil, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, rbac.CodeRoleNameConflict)
}
return role, nil
```

This only applies to the **final return**. Mid-function errors still need `if err != nil` to stop execution.

## Creating Errors

Most of the time, use cases create errors on demand using `errx.New`:

```go
return nil, errx.New(
    "admin not found",
    errx.WithType(errx.T_NotFound),
    errx.WithCode(user.CodeAdminNotFound),
)
```

### Pre-defined Error Variables

Pre-defining package-level error variables is only needed when the **same error is reused multiple times** within a use case. For example, `admin-login` returns the same "incorrect credentials" error for multiple failure reasons (wrong username, inactive account, wrong password):

```go
var errIncorrectCreds = errx.New(
    "username or password is incorrect",
    errx.WithType(errx.T_Validation),
    errx.WithCode(user.CodeIncorrectCreds),
)
```

Use with context to distinguish the cause:

```go
return nil, errx.Wrap(errIncorrectCreds, errx.WithDetails(errx.D{"cause": "password"}))
```

## Error Details

The most appropriate place to add error details is when creating the error with `errx.New` — that's where you have the most context about what went wrong.

When adding whole structs (requests, responses, event bodies) to error details, always **mask them first** using the `mask` package. This ensures sensitive fields (marked with `mask:"true"` tag) like passwords and tokens are not exposed in error logs.
