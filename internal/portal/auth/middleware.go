package auth

import (
	"context"
	"strings"

	"github.com/code19m/errx"
	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/meta"
)

// NewAuthMiddleware returns a Fiber middleware that authenticates requests
// using the Bearer token from the Authorization header.
func NewAuthMiddleware(p Portal) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearerToken(c)
		if token == "" {
			return errx.New(
				"missing or invalid authorization header",
				errx.WithType(errx.T_Authentication),
				errx.WithCode(CodeUnauthorized),
			)
		}

		uc, err := p.Authenticate(c.UserContext(), token)
		if err != nil {
			return errx.Wrap(err)
		}

		ctx := setUserContext(c.UserContext(), uc)
		ctx = context.WithValue(ctx, meta.ActorType, ActorTypeUser)
		ctx = context.WithValue(ctx, meta.ActorID, uc.UserID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// RequirePermission returns a Fiber middleware that checks if the authenticated user
// has at least one of the specified permissions (OR logic).
// For AND logic, chain multiple RequirePermission middleware calls.
func RequirePermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uc := UserContextFromCtx(c.UserContext())
		if uc == nil {
			return errx.New(
				"user context not found",
				errx.WithType(errx.T_Authentication),
				errx.WithCode(CodeUnauthorized),
			)
		}

		for _, perm := range permissions {
			if HasPermission(uc, perm) {
				return c.Next()
			}
		}

		return errx.New(
			"insufficient permissions",
			errx.WithType(errx.T_Forbidden),
			errx.WithCode(CodeForbidden),
			errx.WithDetails(errx.D{"required_permissions": permissions}),
		)
	}
}

func extractBearerToken(c *fiber.Ctx) string {
	header := c.Get("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}
