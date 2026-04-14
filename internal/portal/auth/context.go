package auth

import "context"

type UserContext struct {
	UserID      string
	Username    *string
	Email       *string
	IsVerified  bool
	SessionID   int64
	Permissions []string
}

type contextKey struct{}

// setUserContext stores the UserContext in the given context.
func setUserContext(ctx context.Context, uc *UserContext) context.Context {
	return context.WithValue(ctx, contextKey{}, uc)
}

// UserContextFromCtx retrieves the UserContext from the context, or nil if not present.
func UserContextFromCtx(ctx context.Context) *UserContext {
	uc, _ := ctx.Value(contextKey{}).(*UserContext)
	return uc
}
