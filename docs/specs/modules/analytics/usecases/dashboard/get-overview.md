# Get Dashboard Overview

Returns the authenticated candidate's dashboard overview for the selected date range.

> **type**: user_action

> **operation-id**: `get-dashboard-overview`

> **access**: GET /api/v1/dashboard/overview

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/analytics/usecase/dashboard/getoverview/usecase.go)

## Input

Query parameters:

```json
{
  "range": "7d" // optional, allowed: 7d, 30d, 90d, all. Default: 7d
}
```

## Output

```json
{
  "user": {
    "id": "user_123",
    "full_name": "John Developer", // nullable
    "email": "john@example.com", // nullable
    "avatar_url": null, // nullable
    "target_role": {
      "id": "python",
      "name": "Python Backend"
    }, // nullable
    "experience_level": {
      "id": "mid",
      "name": "Mid-Level"
    } // nullable
  },
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
  },
  "performance": {
    "range": "7d",
    "summary": {
      "average_score": 87, // nullable
      "score_delta_percent": 15,
      "interviews_completed": 8,
      "practice_seconds": 14400
    },
    "points": [
      {
        "date": "2026-04-14",
        "label": "Apr 14",
        "average_score": 87, // nullable
        "interviews_completed": 1,
        "practice_seconds": 3600
      }
    ]
  },
  "topics": {
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
    ],
    "weak": [
      {
        "id": "security",
        "name": "Security",
        "score": 65, // nullable
        "questions_answered": 12,
        "reason": "Lowest average score in the selected range.",
        "recommended_action": "Review authentication, authorization, and SQL injection prevention."
      }
    ],
    "strong": [
      {
        "id": "api_design",
        "name": "API Design",
        "score": 90, // nullable
        "questions_answered": 20,
        "reason": "Consistently high scores across recent sessions."
      }
    ]
  },
  "recent_activity": {
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
    "next_cursor": null // nullable
  },
  "recommendations": {
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
}
```

## Field Rules

- `delta_direction`: allowed values: `up`, `down`, `flat`, `new`
- `trend`: allowed values: `up`, `down`, `flat`, `new`
- `priority`: allowed values: `low`, `medium`, `high`
- `level`: allowed values: `strong`, `stable`, `needs_practice`
- `status`: allowed values: `in_progress`, `completed`, `abandoned`, `scoring`
- `difficulty`: allowed values: `easy`, `medium`, `hard`, `mixed`
- Scores are integers from `0` to `100`, or `null` when unavailable
- Durations are seconds
- Exact timestamps use ISO 8601 UTC
- Chart bucket dates use `YYYY-MM-DD`
- `performance.points` is sorted oldest to newest
- Empty sections are returned as empty arrays, not omitted

## Execute

- Validate input

- Read authenticated user context

- Find candidate profile for authenticated user

- List preferred topics for the candidate profile when it exists

- Find candidate progress summary for authenticated user

- List candidate topic stats for selected range

- List recent interview activity for authenticated user and selected range

- Find current avatar file metadata when it exists

- Build dashboard stats with current values and range deltas

- Build performance chart points sorted oldest to newest

- Build topic performance items, weak topics, and strong topics

- Build next-interview recommendations from weak topics, preferred topics, target role, and experience level

- Return dashboard overview

## Empty State

Return a successful response for candidates without completed interview sessions.

- `stats.total_interviews.value` is `0`
- `stats.average_score.value` is `null`
- `stats.current_streak_days.value` is `0`
- `stats.total_practice_seconds.value` is `0`
- `performance.summary.average_score` is `null`
- `performance.points` is `[]`
- `topics.items` is `[]`
- `topics.weak` is `[]`
- `topics.strong` is `[]`
- `recent_activity.items` is `[]`
- `recommendations.recommended_topics` is `[]`

## Errors

- Return `VALIDATION_FAILED` when `range` is not one of the allowed values
- Return `DASHBOARD_UNAVAILABLE` when a downstream dashboard dependency is temporarily unavailable
