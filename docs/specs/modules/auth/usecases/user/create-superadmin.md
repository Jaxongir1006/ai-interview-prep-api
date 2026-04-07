# Create Superadmin

Creates the initial superadmin account for system bootstrap. This command should be run once during initial setup.

> **type**: manual_command

> **operation-id**: `create-superadmin`

> **usage**: `./app auth create-superadmin`

> **implementation**: [usecase.go](../../../../../../internal/modules/auth/usecase/user/createsuperadmin/usecase.go)

## Execute

- Hash the password

- Start UOW

- Create user

- Assign all `SuperadminPermissions` explicitly (no bypass — real permission checks apply)

- Apply UOW
