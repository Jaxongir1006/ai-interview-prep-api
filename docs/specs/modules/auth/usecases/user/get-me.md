# Get Me

Returns the full current-user snapshot by aggregating auth identity, candidate profile data, preferred topics, progress summary, and avatar metadata.

> **type**: user_action

> **operation-id**: `get-me`

> **access**: GET /api/v1/auth/get-me

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/getme/usecase.go)

## Input

No input required. User identity is provided by auth middleware via the `Authorization` header.

## Output

```json
{
  "user": {
    "id": "string",
    "username": "admin-username", // nullable
    "email": "user@example.com", // nullable
    "phone_number": "+998901234567", // nullable
    "is_verified": false,
    "is_active": true,
    "last_login_at": "2024-01-01T00:00:00Z", // nullable
    "last_active_at": "2024-01-01T00:00:00Z", // nullable
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "oauth_providers": [
      "google",
      "github"
    ]
  },
  "profile": {
    "id": 1,
    "user_id": "string",
    "full_name": "John Doe", // nullable
    "bio": "Preparing for backend interviews", // nullable
    "location": "Tashkent, Uzbekistan", // nullable
    "target_role": "Golang Backend Developer",
    "experience_level": "junior",
    "interview_goal_per_week": 3,
    "preferred_topics": [
      "golang-concurrency",
      "postgres-indexing"
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }, // nullable
  "progress_summary": {
    "current_streak": 2,
    "longest_streak": 5,
    "total_interviews_taken": 8,
    "total_time_spent_seconds": 5400,
    "average_score": 78.5,
    "last_interview_at": "2024-01-01T00:00:00Z" // nullable
  }, // nullable
  "avatar": {
    "file_id": "string",
    "original_filename": "avatar.png",
    "mime_type": "image/png",
    "size_bytes": 123456,
    "download_url": "/api/v1/filevault/download?id=string"
  } // nullable
}
```

## Execute

- Get user ID from authenticated user context

- Find auth user by ID

- List linked OAuth providers for the user

- Find candidate profile for the user

- List preferred topics for the candidate profile when it exists

- Find candidate progress summary when it exists

- Find current avatar file metadata when it exists

- Return aggregated current-user data
