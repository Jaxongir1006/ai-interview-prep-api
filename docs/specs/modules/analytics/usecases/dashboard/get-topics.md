# Get Dashboard Topics

Returns topic-level performance for the authenticated candidate over the selected date range.

> **type**: user_action

> **operation-id**: `get-dashboard-topics`

> **access**: GET /api/v1/dashboard/topics

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/analytics/usecase/dashboard/gettopics/usecase.go)

## Input

Query parameters:

```json
{
  "range": "30d" // optional, allowed: 7d, 30d, 90d, all. Default: 7d
}
```

## Output

```json
{
  "range": "30d",
  "items": [
    {
      "id": "system_design",
      "name": "System Design",
      "score": 72, // nullable
      "questions_answered": 8,
      "correctness_rate": 0.72, // nullable
      "average_time_seconds": 620, // nullable
      "trend": "flat",
      "level": "needs_practice"
    }
  ]
}
```

## Field Rules

- `trend`: allowed values: `up`, `down`, `flat`, `new`
- `level`: allowed values: `strong`, `stable`, `needs_practice`
- `score` is an integer from `0` to `100`, or `null` when unavailable
- `correctness_rate` is a decimal from `0` to `1`, or `null` when unavailable
- `average_time_seconds` is a duration in seconds, or `null` when unavailable
- `items` is sorted by lowest score first, then highest question count

## Execute

- Validate input

- Read authenticated user context

- Read candidate topic stats for authenticated user and selected range

- Read previous comparable range topic stats when available

- Resolve topic display names

- Calculate score, correctness rate, average time, trend, and level for each topic

- Return topic performance items

## Empty State

Return a successful response when no topic performance exists for the selected range.

- `items` is `[]`

## Errors

- Return `VALIDATION_FAILED` when `range` is not one of the allowed values
- Return `DASHBOARD_UNAVAILABLE` when topic performance cannot be calculated temporarily
