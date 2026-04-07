# {Use Case Name}

{One-two sentence summary of what this job does and why it runs on a schedule.}

> **type**: async_task

> **operation-id**: `{operation-id}`

> **implementation**: [usecase.go](../../../../internal/modules/{module}/usecase/{domain}/{operation-id}/usecase.go)

## Task payload

```json
{
  "entity_id": "string", // required
  "action": "string", // required
  "timestamp": "2024-01-15T10:30:00Z" // required, RFC3339
}
```

## Execute

<!--
Describe WHAT the use case does, not HOW.
- Steps should map 1:1 to use case Execute() method logic
- Focus on business actions, not implementation details
- Do NOT include: logging, metrics, query mechanics, infrastructure concerns
- Bad: "Query sessions where X < NOW()" (implementation detail)
- Good: "Delete expired sessions" (business action)
-->

- Validate task payload

- {Business action 1}

- {Business action 2}

## Idempotency

{Brief description of how the job handles reruns - watermark/timestamp tracking, idempotent operations, overlap handling.}

## Schedule

<!-- Include this section only for scheduled tasks. Remove for on-demand tasks. -->

Runs {frequency} via cron pattern: `{cron-pattern}` ({human-readable description}).
