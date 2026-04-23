package auth

import (
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/oauthaccount"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/passwordresettoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/rbac"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/infra/postgres"
	redisinfra "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/infra/redis"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"

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

// GetUserByEmail retrieves a user by email.
// Fails the test if the user is not found.
func GetUserByEmail(t *testing.T, email string) *user.User {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewUserRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	u, err := repo.Get(ctx, user.Filter{Email: &email})
	if err != nil {
		t.Fatalf("GetUserByEmail: failed to get user %q: %v", email, err)
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

func GetEmailVerificationTokensByUserID(
	t *testing.T,
	_ string,
) []emailverificationtoken.EmailVerificationToken {
	t.Helper()

	t.Fatalf("GetEmailVerificationTokensByUserID: Redis-backed tokens cannot be listed by user ID")
	return nil
}

func CurrentEmailVerificationTokenHash(t *testing.T, userID, email string) string {
	t.Helper()

	client := database.GetTestRedis(t)

	ctx, cancel := database.QueryContext()
	defer cancel()

	tokenHash, err := client.Get(ctx, "auth:email_verification:user:"+userID+":"+email).Result()
	if err != nil {
		t.Fatalf("CurrentEmailVerificationTokenHash: failed to get token pointer: %v", err)
	}

	return tokenHash
}

func EmailVerificationTokenPointerExists(t *testing.T, userID, email string) bool {
	t.Helper()

	client := database.GetTestRedis(t)

	ctx, cancel := database.QueryContext()
	defer cancel()

	count, err := client.Exists(ctx, "auth:email_verification:user:"+userID+":"+email).Result()
	if err != nil {
		t.Fatalf("EmailVerificationTokenPointerExists: failed to check token pointer: %v", err)
	}

	return count == 1
}

func GetEmailVerificationTokenByHash(
	t *testing.T,
	tokenHash string,
) *emailverificationtoken.EmailVerificationToken {
	t.Helper()

	client := database.GetTestRedis(t)
	repo := redisinfra.NewEmailVerificationTokenRepo(client)

	ctx, cancel := database.QueryContext()
	defer cancel()

	evt, err := repo.Consume(ctx, tokenHash)
	if err != nil {
		t.Fatalf("GetEmailVerificationTokenByHash: failed to get token: %v", err)
	}

	return evt
}

func ConsumePasswordResetTokenByHash(
	t *testing.T,
	tokenHash string,
) *passwordresettoken.PasswordResetToken {
	t.Helper()

	client := database.GetTestRedis(t)
	repo := redisinfra.NewPasswordResetTokenRepo(client)

	ctx, cancel := database.QueryContext()
	defer cancel()

	resetToken, err := repo.Consume(ctx, tokenHash)
	if err != nil {
		t.Fatalf("ConsumePasswordResetTokenByHash: failed to consume token: %v", err)
	}

	return resetToken
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

// ListOAuthAccountsByUserID returns linked OAuth accounts for the user.
func ListOAuthAccountsByUserID(t *testing.T, userID string) []oauthaccount.OAuthAccount {
	t.Helper()

	db := database.GetTestDB(t)
	repo := postgres.NewOAuthAccountRepo(db)

	ctx, cancel := database.QueryContext()
	defer cancel()

	accounts, err := repo.List(ctx, oauthaccount.Filter{UserID: &userID})
	if err != nil {
		t.Fatalf("ListOAuthAccountsByUserID: failed to list oauth accounts for user %q: %v", userID, err)
	}

	return accounts
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
