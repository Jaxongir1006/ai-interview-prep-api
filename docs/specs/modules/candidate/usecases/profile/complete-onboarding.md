# Complete Onboarding

Completes the authenticated candidate's onboarding by saving selected Interview catalog keys for target role, experience level, and preferred topics.

> **type**: user_action

> **operation-id**: `complete-onboarding`

> **access**: POST /api/v1/me/complete-onboarding

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/candidate/usecase/profile/completeonboarding/usecase.go)

## Input

```json
{
  "target_role": "golang", // required, active Interview target-role key
  "experience_level": "junior", // required, active Interview experience-level key
  "preferred_topics": ["algorithms", "system-design"] // required, min=1, max=10, unique active Interview topic keys
}
```

## Output

```json
{
  "profile": {
    "id": 123,
    "user_id": "user-id",
    "full_name": "John Doe",
    "target_role": "golang",
    "experience_level": "junior",
    "preferred_topics": ["algorithms", "system-design"],
    "onboarding_completed": true,
    "onboarding_completed_at": "2026-04-13T10:00:00Z"
  }
}
```

## Execute

- Validate input

- Read authenticated user context

- Ensure authenticated user is verified

- Validate `target_role`, `experience_level`, and `preferred_topics` against active Interview catalog options

- Find candidate profile by authenticated user ID

- Start UOW

- Update profile `target_role` with selected Interview target-role key

- Update profile `experience_level` with selected Interview experience-level key

- Set profile `onboarding_completed = true`

- Set profile `onboarding_completed_at` to current timestamp

- Replace existing preferred topics with selected Interview topic keys, preserving input order as priority

- Apply UOW

- Return updated profile with preferred topics

## Errors

- Return `CANDIDATE_PROFILE_NOT_FOUND` when the authenticated user does not have a candidate profile
- Return `EMAIL_NOT_VERIFIED` when the authenticated user is not verified
- Return `VALIDATION_FAILED` when `target_role`, `experience_level`, or `preferred_topics` is missing, duplicated, inactive, or not found in the Interview catalog
