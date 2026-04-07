//go:build system

package user_test

import (
	"strings"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	stateauth "github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"

	"github.com/rise-and-shine/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestCreateSuperadmin_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	username := "superadmin"
	password := "securepassword123"

	// WHEN
	output := trigger.ManualCommand(t, "auth", "create-superadmin", "--username", username, "--password", password)

	// THEN
	assert.Contains(t, output, "Superadmin created successfully")

	u := stateauth.GetUserByUsername(t, username)
	assert.Equal(t, username, *u.Username)
	assert.True(t, u.IsActive)

	// Verify password is hashed (not stored in plain text)
	assert.NotEqual(t, password, *u.PasswordHash, "password should not be stored in plain text")
	assert.True(t, strings.HasPrefix(*u.PasswordHash, "$2"), "password hash should be bcrypt format")
	assert.True(t, hasher.Compare(password, *u.PasswordHash), "stored hash should match original password")

	// Verify all superadmin permissions are assigned
	for _, perm := range auth.SuperadminPermissions() {
		assert.True(t, stateauth.HasPermission(t, u.ID, perm),
			"superadmin should have permission: %s", perm)
	}
}

func TestCreateSuperadmin_UsernameConflict(t *testing.T) {
	// GIVEN
	database.Empty(t)
	username := "existingadmin"
	password := "password123"

	stateauth.GivenUsers(t, map[string]any{
		"username": username,
	})

	// WHEN
	output, err := trigger.ManualCommandWithError(
		t,
		"auth",
		"create-superadmin",
		"--username",
		username,
		"--password",
		password,
	)

	// THEN
	assert.Error(t, err, "command should fail for duplicate username")
	assert.Contains(t, strings.ToLower(output), "conflict",
		"error message should indicate conflict: %s", output)
}
