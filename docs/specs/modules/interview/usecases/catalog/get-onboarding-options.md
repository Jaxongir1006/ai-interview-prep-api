# Get Onboarding Options

Returns active interview catalog options used by the frontend to render candidate onboarding choices.

> **type**: user_action

> **operation-id**: `get-onboarding-options`

> **access**: GET /api/v1/interview/get-onboarding-options

> **actor**: user (authenticated)

> **permissions**: -

> **implementation**: [usecase.go](../../../../../../internal/modules/interview/usecase/catalog/getonboardingoptions/usecase.go)

## Input

No input.

## Output

```json
{
  "target_roles": [
    {
      "key": "golang",
      "name": "Go Developer",
      "description": "Backend interviews focused on Go services", // nullable
      "display_order": 10
    }
  ],
  "experience_levels": [
    {
      "key": "junior",
      "name": "Junior",
      "description": "Early-career interview expectations", // nullable
      "display_order": 10
    }
  ],
  "topics": [
    {
      "key": "system-design",
      "name": "System Design",
      "description": "Scalable architecture and distributed systems", // nullable
      "category": "backend", // nullable
      "target_role_keys": ["golang", "java"],
      "display_order": 20
    }
  ]
}
```

## Execute

- Read active target roles ordered by `display_order` and `name`

- Read active experience levels ordered by `display_order` and `name`

- Read active topics ordered by `display_order` and `name`

- Read active target-role topic mappings for active target roles and active topics

- Attach `target_role_keys` to each topic

- Return onboarding option catalog

## Rules

- Return only active target roles, experience levels, topics, and role-topic mappings
- Return stable keys for client submissions and display names for UI labels
- Do not return authorization roles from the Auth module
- `target_role_keys` may be empty when a topic is globally available
- The response is read-only catalog data; it must not create candidate profile state

## Errors

- Return `UNAUTHORIZED` when the request is not authenticated
