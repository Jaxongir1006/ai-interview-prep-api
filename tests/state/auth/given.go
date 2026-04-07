package auth

import (
	"testing"
	"time"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/session"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/infra/postgres"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	"github.com/jaxongir1006/hire-ready-api/pkg/anymap"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"

	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/token"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// GivenUsers creates user records in the database for test setup.
// Each map in data represents a user with the following valid keys:
//   - id: string (default: generated UUID)
//   - username: string (default: "testadmin")
//   - password: string (default: TestPassword1) - will be hashed
//   - is_active: bool (default: true)
//   - last_active_at: *time.Time (default: nil)
//
// Returns the created user entities.
func GivenUsers(t *testing.T, data ...map[string]any) []user.User {
	t.Helper()

	if len(data) == 0 {
		data = []map[string]any{{}}
	}

	db := database.GetTestDB(t)
	repo := postgres.NewUserRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	users := make([]user.User, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenUsers", validUserKeys, d)

		var passwordHash *string
		if v, hasKey := d["password"]; hasKey && v == nil {
			passwordHash = nil
		} else {
			password := anymap.String(d, "password", TestPassword1)
			hashStr, precomputed := GetPrecomputedHash(password)
			if !precomputed {
				var err error
				hashStr, err = hasher.Hash(password)
				if err != nil {
					t.Fatalf("GivenUsers[%d]: failed to hash password: %v", i, err)
				}
			}
			passwordHash = &hashStr
		}

		u := &user.User{
			ID:           anymap.String(d, "id", uuid.NewString()),
			Username:     anymap.StringPtr(d, "username", lo.ToPtr("testadmin")),
			PasswordHash: passwordHash,
			IsActive:     anymap.Bool(d, "is_active", true),
			LastActiveAt: anymap.TimePtr(d, "last_active_at", nil),
		}

		created, err := repo.Create(ctx, u)
		if err != nil {
			t.Fatalf("GivenUsers[%d]: failed to create user: %v", i, err)
		}

		users = append(users, *created)
	}

	return users
}

// GivenSessions creates session records in the database for test setup.
// Each map in data represents a session with the following valid keys:
//   - id: int64 (default: auto-generated)
//   - user_id: string (required)
//   - access_token: string (default: generated opaque token)
//   - access_token_expires_at: time.Time (default: 15 minutes from now)
//   - refresh_token: string (default: generated opaque token)
//   - refresh_token_expires_at: time.Time (default: 30 days from now)
//   - ip_address: string (default: "127.0.0.1")
//   - user_agent: string (default: "test-agent")
//   - last_used_at: time.Time (default: now)
//
// Returns the created session entities.
func GivenSessions(t *testing.T, data ...map[string]any) []session.Session {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenSessions: at least one session data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewSessionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	sessions := make([]session.Session, 0, len(data))
	now := time.Now()

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenSessions", validSessionKeys, d)

		userID := anymap.String(d, "user_id", "")
		if userID == "" {
			t.Fatalf("GivenSessions[%d]: user_id is required", i)
		}

		sess := &session.Session{
			UserID:                userID,
			AccessToken:           anymap.String(d, "access_token", token.NewOpaqueToken()),
			AccessTokenExpiresAt:  anymap.Time(d, "access_token_expires_at", now.Add(15*time.Minute)),
			RefreshToken:          anymap.String(d, "refresh_token", token.NewOpaqueToken()),
			RefreshTokenExpiresAt: anymap.Time(d, "refresh_token_expires_at", now.Add(30*24*time.Hour)),
			IPAddress:             anymap.String(d, "ip_address", "127.0.0.1"),
			UserAgent:             anymap.String(d, "user_agent", "test-agent"),
			LastUsedAt:            anymap.Time(d, "last_used_at", now),
		}

		created, err := repo.Create(ctx, sess)
		if err != nil {
			t.Fatalf("GivenSessions[%d]: failed to create session: %v", i, err)
		}

		sessions = append(sessions, *created)
	}

	return sessions
}

// GivenUserPermissions creates user permission records in the database for test setup.
// Each map in data represents a permission with the following valid keys:
//   - id: int64 (default: auto-generated)
//   - user_id: string (required)
//   - permission: string (required)
//
// Returns the created user permission entities.
func GivenUserPermissions(t *testing.T, data ...map[string]any) []rbac.UserPermission {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenUserPermissions: at least one permission data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewUserPermissionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	permissions := make([]rbac.UserPermission, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenUserPermissions", validUserPermissionKeys, d)

		userID := anymap.String(d, "user_id", "")
		permission := anymap.String(d, "permission", "")

		if userID == "" || permission == "" {
			t.Fatalf("GivenUserPermissions[%d]: user_id and permission are required", i)
		}

		perm := &rbac.UserPermission{
			UserID:     userID,
			Permission: permission,
		}

		created, err := repo.Create(ctx, perm)
		if err != nil {
			t.Fatalf("GivenUserPermissions[%d]: failed to create permission: %v", i, err)
		}

		permissions = append(permissions, *created)
	}

	return permissions
}

// GivenRoles creates role records in the database for test setup.
// Each map in data represents a role with the following valid keys:
//   - id: int64 (default: auto-generated)
//   - name: string (required)
//
// Returns the created role entities.
func GivenRoles(t *testing.T, data ...map[string]any) []rbac.Role {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenRoles: at least one role data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewRoleRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	roles := make([]rbac.Role, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenRoles", validRoleKeys, d)

		name := anymap.String(d, "name", "")
		if name == "" {
			t.Fatalf("GivenRoles[%d]: name is required", i)
		}

		role := &rbac.Role{
			Name: name,
		}

		created, err := repo.Create(ctx, role)
		if err != nil {
			t.Fatalf("GivenRoles[%d]: failed to create role: %v", i, err)
		}

		roles = append(roles, *created)
	}

	return roles
}

// GivenRolePermissions creates role permission records in the database for test setup.
// Each map in data represents a role permission with the following valid keys:
//   - id: int64 (default: auto-generated)
//   - role_id: int64 (required)
//   - permission: string (required)
//
// Returns the created role permission entities.
func GivenRolePermissions(t *testing.T, data ...map[string]any) []rbac.RolePermission {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenRolePermissions: at least one role permission data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewRolePermissionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	perms := make([]rbac.RolePermission, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenRolePermissions", validRolePermissionKeys, d)

		roleID := cast.ToInt64(d["role_id"])
		permission := anymap.String(d, "permission", "")

		if roleID == 0 || permission == "" {
			t.Fatalf("GivenRolePermissions[%d]: role_id and permission are required", i)
		}

		rp := &rbac.RolePermission{
			RoleID:     roleID,
			Permission: permission,
		}

		created, err := repo.Create(ctx, rp)
		if err != nil {
			t.Fatalf("GivenRolePermissions[%d]: failed to create role permission: %v", i, err)
		}

		perms = append(perms, *created)
	}

	return perms
}

// GivenUserRoles creates user role assignment records in the database for test setup.
// Each map in data represents a user role with the following valid keys:
//   - id: int64 (default: auto-generated)
//   - user_id: string (required)
//   - role_id: int64 (required)
//
// Returns the created user role entities.
func GivenUserRoles(t *testing.T, data ...map[string]any) []rbac.UserRole {
	t.Helper()

	if len(data) == 0 {
		t.Fatal("GivenUserRoles: at least one user role data map is required")
	}

	db := database.GetTestDB(t)
	repo := postgres.NewUserRoleRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	userRoles := make([]rbac.UserRole, 0, len(data))

	for i, d := range data {
		anymap.ValidateKeys(t, "GivenUserRoles", validUserRoleKeys, d)

		userID := anymap.String(d, "user_id", "")
		roleID := cast.ToInt64(d["role_id"])

		if userID == "" || roleID == 0 {
			t.Fatalf("GivenUserRoles[%d]: user_id and role_id are required", i)
		}

		ur := &rbac.UserRole{
			UserID: userID,
			RoleID: roleID,
		}

		created, err := repo.Create(ctx, ur)
		if err != nil {
			t.Fatalf("GivenUserRoles[%d]: failed to create user role: %v", i, err)
		}

		userRoles = append(userRoles, *created)
	}

	return userRoles
}

// GivenAuthToken creates a user with the specified permissions and an active session,
// returning the access token for use in Authorization headers.
// Permissions should be portal auth permission constants (e.g., auth.PermissionUserRead).
func GivenAuthToken(t *testing.T, permissions ...string) string {
	t.Helper()

	u := GivenUsers(t, map[string]any{"username": "authuser-" + uuid.NewString()[:8]})[0]
	sessions := GivenSessions(t, map[string]any{"user_id": u.ID})

	for _, perm := range permissions {
		GivenUserPermissions(t, map[string]any{
			"user_id":    u.ID,
			"permission": perm,
		})
	}

	return sessions[0].AccessToken
}

// GivenSuperadminToken creates a superadmin user with all permissions and an active session,
// returning the access token for use in Authorization headers.
func GivenSuperadminToken(t *testing.T) string {
	t.Helper()
	return GivenAuthToken(t, auth.SuperadminPermissions()...)
}

// GivenUserSessions is a convenience function that creates multiple sessions for a user.
// Sessions are created with staggered last_used_at times (1 hour apart) so they have a
// deterministic ordering. Returns sessions ordered by last_used_at ASC (oldest first).
func GivenUserSessions(t *testing.T, userID string, count int) []session.Session {
	t.Helper()

	sessions := make([]session.Session, 0, count)
	now := time.Now()

	for i := range count {
		sess := GivenSessions(t, map[string]any{
			"user_id":      userID,
			"last_used_at": now.Add(time.Duration(i) * time.Hour),
		})
		sessions = append(sessions, sess...)
	}

	return sessions
}
