# Set User Roles

Replaces all role assignments for a user with the provided list. This is a full replacement (set) operation — existing assignments not in the list are removed, new ones are added.

> **type**: user_action

> **operation-id**: `set-user-roles`

> **access**: POST /api/v1/auth/set-user-roles

> **actor**: user

> **permissions**: `auth:access:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/setuserroles/usecase.go)

## Input

```json
{
  "user_id": "string", // required, UUID format
  "role_ids": [1, 2] // required, list of integer role IDs
}
```

## Output

Empty response.

## Execute

- Find user by ID

- Validate all role IDs exist

- Start UOW

- Delete all existing role assignments for the user

- Create new role assignment records

- Record audit log

- Apply UOW
