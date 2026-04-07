# Set User Permissions

Replaces all direct permission assignments for a user with the provided list. This is a full replacement (set) operation — existing permissions not in the list are removed, new ones are added.

> **type**: user_action

> **operation-id**: `set-user-permissions`

> **access**: POST /api/v1/auth/set-user-permissions

> **actor**: user

> **permissions**: `auth:access:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/setuserpermissions/usecase.go)

## Input

```json
{
  "user_id": "string", // required, UUID format
  "permissions": ["string"] // required, list of permission strings
}
```

## Output

Empty response.

## Execute

- Find user by ID

- Start UOW

- Delete all existing direct permissions for the user

- Create new permission records

- Record audit log

- Apply UOW
