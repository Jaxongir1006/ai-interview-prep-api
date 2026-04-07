# {Use Case Name}

{One-two sentence summary of what this command does and when an operator would run it.}

> **type**: manual_command

> **operation-id**: `{operation-id}`

> **usage**: `{app-name} {command} [flags] [arguments]`

> **implementation**: [usecase.go](../../../../internal/modules/{module}/usecase/{domain}/{operation-id}/usecase.go)

## Input

Flags:

- `--flag-name`, `-f`: string, required, description

Arguments:

- `arg_name`: string, required, description

## Execute

<!--
Describe WHAT the use case does, not HOW.
- Steps should map 1:1 to use case Execute() method logic
- Focus on business actions, not implementation details
- Do NOT include: logging, metrics, audit logs, query mechanics, infrastructure concerns
-->

- Validate input arguments and flags

- Check preconditions: {what must be true}

- {Business action}
