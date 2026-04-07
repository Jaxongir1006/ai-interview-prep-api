# {Use Case Name}

{One-two sentence summary of what this subscriber does when the event occurs.}

> **type**: event_subscriber

> **operation-id**: `{operation-id}`

> **event**: `{EventName}`

> **implementation**: [usecase.go](../../../../internal/modules/{module}/usecase/{domain}/{operation-id}/usecase.go)

## Event payload

```json
{
  "entity_id": "string", // required
  "action": "string", // required
  "timestamp": "2024-01-15T10:30:00Z", // required, RFC3339
  "metadata": {} // optional, additional context
}
```

## Handle

<!--
Describe WHAT the use case does, not HOW.
- Steps should map 1:1 to use case Handle() method logic
- Focus on business actions, not implementation details
- Do NOT include: logging, metrics, query mechanics, infrastructure concerns
-->

- Validate event payload

- Check idempotency using {dedup-key}

- {Business action}

- Produce `{OutputEventName}` event (if applicable)

## Idempotency

{Brief description of how duplicate events are handled - event ID tracking, deduplication key, or naturally idempotent operation.}
