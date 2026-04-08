package auth

//nolint:gochecknoglobals // static validation maps for test state
var (
	// validUserKeys defines the allowed keys for user test data.
	// Keys correspond to user.User entity fields.
	validUserKeys = map[string]bool{
		"id":             true,
		"username":       true,
		"email":          true,
		"phone_number":   true,
		"password":       true, // will be hashed to password_hash
		"is_verified":    true,
		"is_active":      true,
		"last_login_at":  true,
		"last_active_at": true,
	}

	// validSessionKeys defines the allowed keys for session test data.
	// Keys correspond to session.Session entity fields.
	validSessionKeys = map[string]bool{
		"id":                       true,
		"user_id":                  true,
		"access_token":             true,
		"access_token_expires_at":  true,
		"refresh_token":            true,
		"refresh_token_expires_at": true,
		"ip_address":               true,
		"user_agent":               true,
		"last_used_at":             true,
	}

	// validUserPermissionKeys defines the allowed keys for user permission test data.
	// Keys correspond to rbac.UserPermission entity fields.
	validUserPermissionKeys = map[string]bool{
		"id":         true,
		"user_id":    true,
		"permission": true,
	}

	// validRoleKeys defines the allowed keys for role test data.
	// Keys correspond to rbac.Role entity fields.
	validRoleKeys = map[string]bool{
		"id":   true,
		"name": true,
	}

	// validUserRoleKeys defines the allowed keys for user role test data.
	// Keys correspond to rbac.UserRole entity fields.
	validUserRoleKeys = map[string]bool{
		"id":      true,
		"user_id": true,
		"role_id": true,
	}

	// validRolePermissionKeys defines the allowed keys for role permission test data.
	// Keys correspond to rbac.RolePermission entity fields.
	validRolePermissionKeys = map[string]bool{
		"id":         true,
		"role_id":    true,
		"permission": true,
	}
)
