//go:build system

package user_test

import (
	"net/http"
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/auth"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
	"github.com/Jaxongir1006/ai-interview-prep-api/tests/system/trigger"
)

func TestGetUsers_Success(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	auth.GivenUsers(t,
		map[string]any{"username": "alice"},
		map[string]any{"username": "bob"},
		map[string]any{"username": "charlie", "is_active": false},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-users").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.Value("page_number").Number().IsEqual(1)
	obj.Value("page_size").Number().IsEqual(20)
	obj.Value("count").Number().IsEqual(4)
	obj.Value("content").Array().Length().IsEqual(4) // 3 given + 1 superadmin from token
}

func TestGetUsers_FilterByIsActive(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)
	auth.GivenUsers(t,
		map[string]any{"username": "active1"},
		map[string]any{"username": "active2"},
		map[string]any{"username": "inactive", "is_active": false},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-users").
		WithQuery("is_active", true).
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	resp.JSON().Object().Value("content").Array().Length().IsEqual(3) // 2 given active + 1 superadmin from token
}

func TestGetUsers_WithRolesAndPermissions(t *testing.T) {
	// GIVEN
	database.Empty(t)
	token := auth.GivenSuperadminToken(t)

	users := auth.GivenUsers(t,
		map[string]any{"username": "alice"},
		map[string]any{"username": "bob"},
	)
	alice := users[0]
	bob := users[1]

	roles := auth.GivenRoles(t,
		map[string]any{"name": "admin"},
		map[string]any{"name": "editor"},
	)
	adminRole := roles[0]
	editorRole := roles[1]

	auth.GivenUserRoles(t,
		map[string]any{"user_id": alice.ID, "role_id": adminRole.ID},
		map[string]any{"user_id": alice.ID, "role_id": editorRole.ID},
		map[string]any{"user_id": bob.ID, "role_id": editorRole.ID},
	)

	auth.GivenUserPermissions(t,
		map[string]any{"user_id": alice.ID, "permission": "auth:user:read"},
		map[string]any{"user_id": bob.ID, "permission": "auth:role:read"},
		map[string]any{"user_id": bob.ID, "permission": "auth:role:write"},
	)

	// WHEN
	resp := trigger.UserAction(t).GET("/api/v1/auth/get-users").
		WithHeader("Authorization", "Bearer "+token).
		Expect()

	// THEN
	resp.Status(http.StatusOK)
	arr := resp.JSON().Object().Value("content").Array()
	arr.Length().IsEqual(3) // alice + bob + superadmin from token

	// Find alice and bob in response
	for _, item := range arr.Iter() {
		obj := item.Object()
		username := obj.Value("username").String().Raw()

		switch username {
		case "alice":
			obj.Value("roles").Array().Length().IsEqual(2)
			obj.Value("roles").Array().ContainsAny("admin", "editor")
			obj.Value("direct_permissions").Array().Length().IsEqual(1)
			obj.Value("direct_permissions").Array().ContainsAny("auth:user:read")
		case "bob":
			obj.Value("roles").Array().Length().IsEqual(1)
			obj.Value("roles").Array().ContainsAny("editor")
			obj.Value("direct_permissions").Array().Length().IsEqual(2)
			obj.Value("direct_permissions").Array().ContainsAny("auth:role:read", "auth:role:write")
		}
	}
}

func TestGetUsers_AuthFailures(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing authorization header",
			setup:      func(_ *testing.T) string { return "" },
			wantStatus: http.StatusUnauthorized,
			wantErr:    "UNAUTHORIZED",
		},
		{
			name: "insufficient permissions",
			setup: func(t *testing.T) string {
				return auth.GivenAuthToken(t)
			},
			wantStatus: http.StatusForbidden,
			wantErr:    "FORBIDDEN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			database.Empty(t)
			token := tc.setup(t)

			// WHEN
			req := trigger.UserAction(t).GET("/api/v1/auth/get-users")
			if token != "" {
				req = req.WithHeader("Authorization", "Bearer "+token)
			}
			resp := req.Expect()

			// THEN
			resp.Status(tc.wantStatus)
			resp.JSON().Object().Value("error").Object().
				Value("code").String().IsEqual(tc.wantErr)
			resp.JSON().Object().Value("error").Object().
				Value("message").String().NotContains("[untranslated]")
		})
	}
}
