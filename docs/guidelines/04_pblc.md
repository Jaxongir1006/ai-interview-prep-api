# PBLC

Packaged Business Logic Components — a helper layer for use cases.

## Why PBLC Exists

Our use cases are command-based — each use case is a separate struct with its own `Execute` method. This means there is no shared struct where common logic can live. You can't simply extract a method on the use case struct and reuse it across multiple use cases.

PBLC solves this by providing a dedicated layer where shared or complex business logic lives, callable from any use case.

## When to Use PBLC

- **Deduplication** — logic that repeats across multiple use cases
- **Complex business logic** — state machines, strategy patterns, or any logic worth encapsulating clearly on its own

## Design Freedom

PBLC has no prescribed structure. Choose whatever design fits the problem:

- **Simple functions** — for deduplication, just write functions that accept dependencies as input parameters, do the logic, and return results. No structs or interfaces needed.
- **Structs with interfaces** — for complex components that benefit from encapsulation (e.g., a payment processor, a permission checker)
- **OOP patterns** — State, Strategy, Builder, etc. when the business logic genuinely calls for it

## PBLC Rules

- **Called only from UCs** — never from controllers or other layers
- **Validate inputs strictly** — PBLC doesn't know its caller
- **Return error codes** — use case layer assigns error types
