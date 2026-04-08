# Register

Registers a public platform user with email and password, then initializes a minimal candidate profile that can be completed later.

> **type**: user_action

> **operation-id**: `register`

> **access**: POST /api/v1/auth/register

> **actor**: user (unauthenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/register/usecase.go)

## Input

```json
{
  "email": "user@example.com", // required, email format, max=255
  "password": "string", // required, min=8, max=72
  "full_name": "John Doe" // required, min=1, max=255
}
```

## Output

```json
{
  "user": {
    "id": "string",
    "email": "user@example.com",
    "phone_number": null,
    "is_verified": false,
    "is_active": true,
    "last_login_at": null,
    "last_active_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "profile": {
    "id": 1,
    "user_id": "string",
    "full_name": "John Doe",
    "bio": null,
    "location": null,
    "target_role": null,
    "experience_level": null,
    "interview_goal_per_week": 3,
    "preferred_topics": [],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## Execute

- Validate input

- Check that no user already exists with the same email

- Start UOW

- Create auth user with email-based credentials

- Create minimal candidate profile for the new user using the provided full name

- Record audit log

- Apply UOW

- Return created user and profile data
