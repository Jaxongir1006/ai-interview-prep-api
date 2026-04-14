# Get Dashboard Stats

Returns the authenticated candidate's dashboard summary cards for the selected date range.

> **type**: user_action

> **operation-id**: `get-dashboard-stats`

> **access**: GET /api/v1/dashboard/stats

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/analytics/usecase/dashboard/getstats/usecase.go)

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
  "stats": {
    "total_interviews": {
      "value": 48,
      "delta_percent": 12,
      "delta_direction": "up"
    },
    "average_score": {
      "value": 87, // nullable
      "delta_percent": 5,
      "delta_direction": "up"
    },
    "current_streak_days": {
      "value": 7,
      "is_record": true
    },
    "total_practice_seconds": {
      "value": 86400,
      "delta_percent": 18,
      "delta_direction": "up"
    }
  }
}
```

## Field Rules

- `delta_direction`: allowed values: `up`, `down`, `flat`, `new`
- Scores are integers from `0` to `100`, or `null` when unavailable
- Durations are seconds
- `delta_percent` is `0` when no comparison data exists

## Execute

- Validate input

- Read authenticated user context

- Find candidate progress summary for authenticated user

- Calculate current range metrics from completed interview activity

- Calculate previous comparable range metrics when available

- Build total interviews, average score, streak, and total practice time stats

- Return dashboard stats

## Empty State

Return a successful response when the candidate has no completed interview sessions.

- `stats.total_interviews.value` is `0`
- `stats.average_score.value` is `null`
- `stats.current_streak_days.value` is `0`
- `stats.current_streak_days.is_record` is `false`
- `stats.total_practice_seconds.value` is `0`
- All `delta_percent` values are `0`
- All `delta_direction` values are `flat`

## Errors

- Return `VALIDATION_FAILED` when `range` is not one of the allowed values
- Return `DASHBOARD_UNAVAILABLE` when dashboard stats cannot be calculated temporarily
