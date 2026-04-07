package auth

import portalaudit "go-enterprise-blueprint/internal/portal/audit"

const (
	PermissionUserRead      = "auth:user:read"
	PermissionUserManage    = "auth:user:manage"
	PermissionRoleRead      = "auth:role:read"
	PermissionRoleManage    = "auth:role:manage"
	PermissionAccessRead    = "auth:access:read"
	PermissionAccessManage  = "auth:access:manage"
	PermissionSessionRead   = "auth:session:read"
	PermissionSessionManage = "auth:session:manage"

	PermissionTaskmillView   = "taskmill:view"
	PermissionTaskmillManage = "taskmill:manage"

	PermissionAlertView   = "alert:view"
	PermissionAlertManage = "alert:manage"
)

// SuperadminPermissions returns the explicit list of permissions assigned to the superadmin user.
// When new permissions are added, they must be included here for superadmin access.
func SuperadminPermissions() []string {
	return []string{
		PermissionUserRead,
		PermissionUserManage,
		PermissionRoleRead,
		PermissionRoleManage,
		PermissionAccessRead,
		PermissionAccessManage,
		PermissionSessionRead,
		PermissionSessionManage,
		PermissionTaskmillView,
		PermissionTaskmillManage,
		PermissionAlertView,
		PermissionAlertManage,
		portalaudit.PermissionActionLogRead,
		portalaudit.PermissionStatusChangeLogRead,
	}
}
