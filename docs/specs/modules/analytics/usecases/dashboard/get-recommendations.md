# Get Dashboard Recommendations

Returns recommended practice topics and next-interview defaults for the authenticated candidate.

> **type**: user_action

> **operation-id**: `get-dashboard-recommendations`

> **access**: GET /api/v1/dashboard/recommendations

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/analytics/usecase/dashboard/getrecommendations/usecase.go)

## Input

No input required. User identity is provided by auth middleware via the `Authorization` header.

## Output

```json
{
  "recommended_topics": [
    {
      "id": "security",
      "name": "Security",
      "priority": "high",
      "reason": "Your Security score is below your overall average."
    }
  ],
  "next_interview": {
    "target_role": {
      "id": "python",
      "name": "Python Backend"
    }, // nullable
    "experience_level": {
      "id": "mid",
      "name": "Mid-Level"
    }, // nullable
    "topics": [
      {
        "id": "security",
        "name": "Security"
      }
    ],
    "difficulty": "medium",
    "question_count": 5,
    "estimated_duration_seconds": 3600
  }
}
```

## Field Rules

- `priority`: allowed values: `low`, `medium`, `high`
- `difficulty`: allowed values: `easy`, `medium`, `hard`, `mixed`
- `recommended_topics` is sorted by priority first, then lowest score
- `next_interview.question_count` defaults to `5`
- `next_interview.estimated_duration_seconds` defaults to `3600`
- `next_interview.difficulty` defaults to `mixed` when the candidate has no enough history

## Execute

- Read authenticated user context

- Find candidate profile for authenticated user

- List preferred topics for the candidate profile when it exists

- Read candidate topic stats

- Find weak topics from lowest scores and low practice volume

- Build recommended topics with priority and reason

- Build next-interview defaults from target role, experience level, preferred topics, and weak topics

- Return dashboard recommendations

## Empty State

Return a successful response when no interview history exists.

- `recommended_topics` is `[]`
- `next_interview.target_role` is copied from candidate profile when available
- `next_interview.experience_level` is copied from candidate profile when available
- `next_interview.topics` is copied from preferred topics when available, otherwise `[]`
- `next_interview.difficulty` is `mixed`
- `next_interview.question_count` is `5`
- `next_interview.estimated_duration_seconds` is `3600`

## Errors

- Return `DASHBOARD_UNAVAILABLE` when recommendations cannot be calculated temporarily
