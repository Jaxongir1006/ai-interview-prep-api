# Get Performance Trend

Returns chart points for the authenticated candidate's performance over the selected date range.

> **type**: user_action

> **operation-id**: `get-performance-trend`

> **access**: GET /api/v1/dashboard/performance-trend

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/analytics/usecase/dashboard/getperformancetrend/usecase.go)

## Input

Query parameters:

```json
{
  "range": "30d", // optional, allowed: 7d, 30d, 90d, all. Default: 7d
  "topic_id": "algorithms" // optional, stable topic identifier
}
```

## Output

```json
{
  "range": "30d",
  "topic": {
    "id": "algorithms",
    "name": "Algorithms"
  }, // nullable
  "points": [
    {
      "date": "2026-04-14",
      "label": "Apr 14",
      "average_score": 87, // nullable
      "interviews_completed": 1,
      "practice_seconds": 3600
    }
  ]
}
```

## Field Rules

- `topic` is `null` when no `topic_id` filter is provided
- `date` uses `YYYY-MM-DD`
- `label` is display-friendly and stable for the selected bucket
- `average_score` is an integer from `0` to `100`, or `null` when unavailable
- `practice_seconds` is a duration in seconds
- `points` is sorted oldest to newest

## Execute

- Validate input

- Read authenticated user context

- Resolve topic display name when `topic_id` is provided

- Read completed interview activity for authenticated user and selected range

- Filter activity by topic when `topic_id` is provided

- Group activity into date buckets

- Calculate average score, completed interview count, and practice seconds for each bucket

- Return performance trend points sorted oldest to newest

## Empty State

Return a successful response when no completed interview activity exists for the selected range.

- `points` is `[]`
- `topic` is returned when a valid `topic_id` filter was provided

## Errors

- Return `VALIDATION_FAILED` when `range` is not one of the allowed values
- Return `VALIDATION_FAILED` when `topic_id` is malformed
- Return `DASHBOARD_UNAVAILABLE` when performance trend data cannot be calculated temporarily
