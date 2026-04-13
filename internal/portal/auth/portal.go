package auth

import (
	"context"
	"slices"
)

const (
	ModuleName = "auth"

	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeSessionExpired   = "SESSION_EXPIRED"
	CodeUserInactive     = "USER_INACTIVE"
	CodeEmailNotVerified = "EMAIL_NOT_VERIFIED"

	ActorTypeUser = "user"
)

type Portal interface {
	// Authenticate authenticates the user and returns a UserContext if successful.
	Authenticate(ctx context.Context, accessToken string) (*UserContext, error)
}

// HasPermission checks if the user has the given permission.
func HasPermission(uc *UserContext, perm string) bool {
	return slices.Contains(uc.Permissions, perm)
}

// MustUserContext retrieves the UserContext from the context.
// Panics if the UserContext is not present.
func MustUserContext(ctx context.Context) *UserContext {
	uc := UserContextFromCtx(ctx)
	if uc == nil {
		panic("auth: UserContext not found in context")
	}
	return uc
}
