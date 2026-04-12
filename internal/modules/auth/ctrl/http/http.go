package http

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/forward"
)

type Controller struct {
	usecaseContainer *usecase.Container
	portalContainer  *portal.Container
	authPortal       auth.Portal
}

func NewContoller(
	usecaseContainer *usecase.Container,
	portalContainer *portal.Container,
	authPortal auth.Portal,
	httpServer *server.HTTPServer,
) *Controller {
	ctrl := &Controller{
		usecaseContainer: usecaseContainer,
		portalContainer:  portalContainer,
		authPortal:       authPortal,
	}

	httpServer.RegisterRouter(ctrl.initRoutes)
	return ctrl
}

func (c *Controller) initRoutes(r fiber.Router) {
	v1 := r.Group("/api/v1/auth")

	// Unauthenticated
	v1.Post("/admin-login", forward.ToUserAction(c.usecaseContainer.AdminLogin()))
	v1.Post("/register", forward.ToUserAction(c.usecaseContainer.Register()))
	v1.Post("/login", forward.ToUserAction(c.usecaseContainer.Login()))
	v1.Post("/google-oauth-login", forward.ToUserAction(c.usecaseContainer.GoogleOAuthLogin()))
	v1.Post("/github-oauth-login", forward.ToUserAction(c.usecaseContainer.GitHubOAuthLogin()))
	v1.Post("/refresh-token", forward.ToUserAction(c.usecaseContainer.RefreshToken()))
	v1.Post("/verify-email", forward.ToUserAction(c.usecaseContainer.VerifyEmail()))
	v1.Post("/resend-verification-email", forward.ToUserAction(c.usecaseContainer.ResendVerificationEmail()))

	// Authenticated (no specific permission required)
	v1Auth := v1.Group("", auth.NewAuthMiddleware(c.authPortal))
	v1Auth.Get("/get-me", forward.ToUserAction(c.usecaseContainer.GetMe()))
	v1Auth.Post("/logout", forward.ToUserAction(c.usecaseContainer.Logout()))
	v1Auth.Post("/change-my-password", forward.ToUserAction(c.usecaseContainer.ChangeMyPassword()))
	v1Auth.Get("/get-my-sessions", forward.ToUserAction(c.usecaseContainer.GetMySessions()))
	v1Auth.Post("/delete-my-session", forward.ToUserAction(c.usecaseContainer.DeleteMySession()))

	// User management
	userRead := auth.RequirePermission(auth.PermissionUserRead)
	userManage := auth.RequirePermission(auth.PermissionUserManage)
	v1Auth.Get("/get-auth-stats", userRead, forward.ToUserAction(c.usecaseContainer.GetAuthStats()))
	v1Auth.Get("/get-users", userRead, forward.ToUserAction(c.usecaseContainer.GetUsers()))
	v1Auth.Post("/create-user", userManage, forward.ToUserAction(c.usecaseContainer.CreateUser()))
	v1Auth.Post("/update-user", userManage, forward.ToUserAction(c.usecaseContainer.UpdateUser()))
	v1Auth.Post("/disable-user", userManage, forward.ToUserAction(c.usecaseContainer.DisableUser()))
	v1Auth.Post("/enable-user", userManage, forward.ToUserAction(c.usecaseContainer.EnableUser()))

	// Role management
	roleRead := auth.RequirePermission(auth.PermissionRoleRead)
	roleManage := auth.RequirePermission(auth.PermissionRoleManage)
	v1Auth.Post("/create-role", roleManage, forward.ToUserAction(c.usecaseContainer.CreateRole()))
	v1Auth.Post("/update-role", roleManage, forward.ToUserAction(c.usecaseContainer.UpdateRole()))
	v1Auth.Post("/delete-role", roleManage, forward.ToUserAction(c.usecaseContainer.DeleteRole()))
	v1Auth.Get("/get-roles", roleRead, forward.ToUserAction(c.usecaseContainer.GetRoles()))
	v1Auth.Post("/set-role-permissions", roleManage,
		forward.ToUserAction(c.usecaseContainer.SetRolePermissions()))
	v1Auth.Get("/get-role-permissions", roleRead,
		forward.ToUserAction(c.usecaseContainer.GetRolePermissions()))

	// Access management
	accessRead := auth.RequirePermission(auth.PermissionAccessRead)
	accessManage := auth.RequirePermission(auth.PermissionAccessManage)
	v1Auth.Post("/set-user-roles", accessManage,
		forward.ToUserAction(c.usecaseContainer.SetUserRoles()))
	v1Auth.Get("/get-user-roles", accessRead,
		forward.ToUserAction(c.usecaseContainer.GetUserRoles()))
	v1Auth.Post("/set-user-permissions", accessManage,
		forward.ToUserAction(c.usecaseContainer.SetUserPermissions()))
	v1Auth.Get("/get-user-permissions", accessRead,
		forward.ToUserAction(c.usecaseContainer.GetUserPermissions()))

	// Session management
	sessionRead := auth.RequirePermission(auth.PermissionSessionRead)
	sessionManage := auth.RequirePermission(auth.PermissionSessionManage)
	v1Auth.Get("/get-user-sessions", sessionRead,
		forward.ToUserAction(c.usecaseContainer.GetUserSessions()))
	v1Auth.Post("/delete-session", sessionManage,
		forward.ToUserAction(c.usecaseContainer.DeleteSession()))
	v1Auth.Post("/delete-user-sessions", sessionManage,
		forward.ToUserAction(c.usecaseContainer.DeleteUserSessions()))
}
