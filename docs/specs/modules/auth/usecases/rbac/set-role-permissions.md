# Set Role Permissions

Replaces all permissions for a role with the provided list. This is a full replacement (set) operation — existing permissions not in the list are removed, new ones are added.

> **type**: user_action

> **operation-id**: `set-role-permissions`

> **access**: POST /api/v1/auth/set-role-permissions

> **actor**: user

> **permissions**: `auth:role:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/setrolepermissions/usecase.go)

## Input

```json
{
  "role_id": 1, // required, integer
  "permissions": ["string"] // required, list of permission strings
}
```

## Output

Empty response.

## Execute

- Find role by ID

- Start UOW

- Delete all existing permissions for the role

- Create new permission records

- Record audit log

- Apply UOW
