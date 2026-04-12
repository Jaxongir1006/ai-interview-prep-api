package domain

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/mail"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/oauthaccount"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/rbac"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
)

// Container holds domain interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct {
	userRepo                   user.Repo
	emailVerificationTokenRepo emailverificationtoken.Repo
	mailSender                 mail.Sender
	oauthAccountRepo           oauthaccount.Repo
	sessionRepo                session.Repo
	roleRepo                   rbac.RoleRepo
	rolePermissionRepo         rbac.RolePermissionRepo
	userRoleRepo               rbac.UserRoleRepo
	userPermissionRepo         rbac.UserPermissionRepo
	uowFactory                 uow.Factory
}

func NewContainer(
	userRepo user.Repo,
	emailVerificationTokenRepo emailverificationtoken.Repo,
	mailSender mail.Sender,
	oauthAccountRepo oauthaccount.Repo,
	sessionRepo session.Repo,
	roleRepo rbac.RoleRepo,
	rolePermissionRepo rbac.RolePermissionRepo,
	userRoleRepo rbac.UserRoleRepo,
	userPermissionRepo rbac.UserPermissionRepo,
	uowFactory uow.Factory,
) *Container {
	return &Container{
		userRepo,
		emailVerificationTokenRepo,
		mailSender,
		oauthAccountRepo,
		sessionRepo,
		roleRepo,
		rolePermissionRepo,
		userRoleRepo,
		userPermissionRepo,
		uowFactory,
	}
}

func (c *Container) UserRepo() user.Repo {
	return c.userRepo
}

func (c *Container) EmailVerificationTokenRepo() emailverificationtoken.Repo {
	return c.emailVerificationTokenRepo
}

func (c *Container) MailSender() mail.Sender {
	return c.mailSender
}

func (c *Container) OAuthAccountRepo() oauthaccount.Repo {
	return c.oauthAccountRepo
}

func (c *Container) SessionRepo() session.Repo {
	return c.sessionRepo
}

func (c *Container) RoleRepo() rbac.RoleRepo {
	return c.roleRepo
}

func (c *Container) RolePermissionRepo() rbac.RolePermissionRepo {
	return c.rolePermissionRepo
}

func (c *Container) UserRoleRepo() rbac.UserRoleRepo {
	return c.userRoleRepo
}

func (c *Container) UserPermissionRepo() rbac.UserPermissionRepo {
	return c.userPermissionRepo
}

func (c *Container) UOWFactory() uow.Factory {
	return c.uowFactory
}
