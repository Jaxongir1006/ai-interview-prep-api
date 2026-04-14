# Get Recent Activity

Returns recent interview sessions for the authenticated candidate.

> **type**: user_action

> **operation-id**: `get-recent-activity`

> **access**: GET /api/v1/dashboard/recent-activity

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/analytics/usecase/dashboard/getrecentactivity/usecase.go)

## Input

Query parameters:

```json
{
  "limit": 10, // optional, min=1, max=50. Default: 10
  "cursor": "session_123" // optional, opaque pagination cursor
}
```

## Output

```json
{
  "items": [
    {
      "session_id": "session_124",
      "title": "System Design Session",
      "status": "completed",
      "score": 78, // nullable
      "started_at": "2026-04-13T14:00:00Z",
      "completed_at": "2026-04-13T15:00:00Z", // nullable
      "duration_seconds": 3600,
      "question_count": 5,
      "answered_count": 5,
      "topics": [
        {
          "id": "system_design",
          "name": "System Design"
        }
      ]
    }
  ],
  "next_cursor": "session_120" // nullable
}
```

## Field Rules

- `status`: allowed values: `in_progress`, `completed`, `abandoned`, `scoring`
- `score` is an integer from `0` to `100`, or `null` when unavailable
- `started_at` and `completed_at` use ISO 8601 UTC timestamps
- `duration_seconds` is a duration in seconds
- `items` is sorted newest to oldest by `started_at`
- `next_cursor` is `null` when there are no more results

## Execute

- Validate input

- Read authenticated user context

- Decode pagination cursor when provided

- List recent interview sessions for authenticated user

- List topics for returned interview sessions

- Build next cursor when more results exist

- Return recent activity items

## Empty State

Return a successful response when no interview sessions exist.

- `items` is `[]`
- `next_cursor` is `null`

## Errors

- Return `VALIDATION_FAILED` when `limit` is outside the allowed range
- Return `VALIDATION_FAILED` when `cursor` is invalid or expired
- Return `DASHBOARD_UNAVAILABLE` when recent activity cannot be loaded temporarily
