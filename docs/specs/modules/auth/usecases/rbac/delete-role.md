# Delete Role

Deletes a role and all its associated permissions. Fails if users are still assigned to the role.

> **type**: user_action

> **operation-id**: `delete-role`

> **access**: POST /api/v1/auth/delete-role

> **actor**: user

> **permissions**: `auth:role:manage`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/rbac/deleterole/usecase.go)

## Input

```json
{
  "id": 1 // required, integer
}
```

## Output

Empty response.

## Execute

- Find role by ID

- Check role has no assigned users

- Start UOW

- Delete all role permissions

- Delete role

- Record audit log

- Apply UOW
