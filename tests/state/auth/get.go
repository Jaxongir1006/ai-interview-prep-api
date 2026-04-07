package auth

import (
	"testing"

	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/rbac"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/session"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/infra/postgres"
	"github.com/jaxongir1006/hire-ready-api/tests/state/database"

	"github.com/rise-and-shine/pkg/sorter"
)

// GetUserByID retrieves a user by ID.
// Fails the test if the user is not found.
func GetUserByID(t *testing.T, id string) *user.User {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	u, err := repo.Get(ctx, user.Filter{ID: &id})
	if err != nil {
		t.Fatalf("GetUserByID: failed to get user %q: %v", id, err)
	}

	return u
}

// GetUserByUsername retrieves a user by username.
// Fails the test if the user is not found.
func GetUserByUsername(t *testing.T, username string) *user.User {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	u, err := repo.Get(ctx, user.Filter{Username: &username})
	if err != nil {
		t.Fatalf("GetUserByUsername: failed to get user %q: %v", username, err)
	}

	return u
}

// UserExists checks if a user with the given username exists.
func UserExists(t *testing.T, username string) bool {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	exists, err := repo.Exists(ctx, user.Filter{Username: &username})
	if err != nil {
		t.Fatalf("UserExists: failed to check user %q: %v", username, err)
	}

	return exists
}

// GetSessionsByUserID retrieves all sessions for a user, ordered by last_used_at ASC.
func GetSessionsByUserID(t *testing.T, userID string) []*session.Session {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewSessionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	asc := sorter.Asc

	sessions, err := repo.List(ctx, session.Filter{
		UserID:            &userID,
		OrderByLastUsedAt: &asc,
	})
	if err != nil {
		t.Fatalf("GetSessionsByUserID: failed to get sessions for user %q: %v", userID, err)
	}

	// Convert to pointer slice
	result := make([]*session.Session, len(sessions))
	for i := range sessions {
		result[i] = &sessions[i]
	}

	return result
}

// SessionExists checks if a session with the given ID exists.
func SessionExists(t *testing.T, sessionID int64) bool {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewSessionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	exists, err := repo.Exists(ctx, session.Filter{ID: &sessionID})
	if err != nil {
		t.Fatalf("SessionExists: failed to check session %d: %v", sessionID, err)
	}

	return exists
}

// SessionCount returns the number of sessions for a user.
func SessionCount(t *testing.T, userID string) int {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewSessionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := repo.Count(ctx, session.Filter{
		UserID: &userID,
	})
	if err != nil {
		t.Fatalf("SessionCount: failed to count sessions for user %q: %v", userID, err)
	}

	return count
}

// GetRoleByID retrieves a role by ID.
// Fails the test if the role is not found.
func GetRoleByID(t *testing.T, id int64) *rbac.Role {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewRoleRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	role, err := repo.Get(ctx, rbac.RoleFilter{ID: &id})
	if err != nil {
		t.Fatalf("GetRoleByID: failed to get role %d: %v", id, err)
	}

	return role
}

// RoleExists checks if a role with the given ID exists.
func RoleExists(t *testing.T, id int64) bool {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewRoleRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	exists, err := repo.Exists(ctx, rbac.RoleFilter{ID: &id})
	if err != nil {
		t.Fatalf("RoleExists: failed to check role %d: %v", id, err)
	}

	return exists
}

// GetRolePermissions returns all permissions assigned to a role.
func GetRolePermissions(t *testing.T, roleID int64) []rbac.RolePermission {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewRolePermissionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	perms, err := repo.List(ctx, rbac.RolePermissionFilter{RoleID: &roleID})
	if err != nil {
		t.Fatalf("GetRolePermissions: failed to get role permissions for role %d: %v", roleID, err)
	}

	return perms
}

// GetUserRoles returns all role assignments for a user.
func GetUserRoles(t *testing.T, userID string) []rbac.UserRole {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserRoleRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	roles, err := repo.List(ctx, rbac.UserRoleFilter{UserID: &userID})
	if err != nil {
		t.Fatalf("GetUserRoles: failed to get user roles for user %q: %v", userID, err)
	}

	return roles
}

// GetUserPermissions returns all direct permissions for a user.
func GetUserPermissions(t *testing.T, userID string) []rbac.UserPermission {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserPermissionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	perms, err := repo.List(ctx, rbac.UserPermissionFilter{UserID: &userID})
	if err != nil {
		t.Fatalf("GetUserPermissions: failed to get user permissions for user %q: %v", userID, err)
	}

	return perms
}

// HasPermission checks if a user has a specific direct permission.
func HasPermission(t *testing.T, userID, permission string) bool {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserPermissionRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	permissions, err := repo.List(ctx, rbac.UserPermissionFilter{
		UserID: &userID,
	})
	if err != nil {
		t.Fatalf("HasPermission: failed to get permissions for user %q: %v", userID, err)
	}

	for _, p := range permissions {
		if p.Permission == permission {
			return true
		}
	}

	return false
}
