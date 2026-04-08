package usecase

import (
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/createrole"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/deleterole"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getrolepermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getroles"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getuserpermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/getuserroles"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/setrolepermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/setuserpermissions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/setuserroles"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/rbac/updaterole"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/cleanexpiredsessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/deletemysession"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/deletesession"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/deleteusersessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/getmysessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/session/getusersessions"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/adminlogin"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/changemypassword"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/createsuperadmin"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/createuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/disableuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/enableuser"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/getauthstats"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/getme"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/getusers"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/githuboauthlogin"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/googleoauthlogin"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/login"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/logout"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/refreshtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/register"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/usecase/user/updateuser"
)

//nolint:dupl // container fields intentionally mirror constructor args for DI clarity
type Container struct {
	// Self-service
	adminLogin       adminlogin.UseCase
	register         register.UseCase
	login            login.UseCase
	getMe            getme.UseCase
	googleOAuthLogin googleoauthlogin.UseCase
	gitHubOAuthLogin githuboauthlogin.UseCase
	refreshToken     refreshtoken.UseCase
	logout           logout.UseCase
	changeMyPassword changemypassword.UseCase

	// User management
	createSuperadmin createsuperadmin.UseCase
	getUsers         getusers.UseCase
	getAuthStats     getauthstats.UseCase
	createUser       createuser.UseCase
	updateUser       updateuser.UseCase
	disableUser      disableuser.UseCase
	enableUser       enableuser.UseCase

	// Role management
	createRole         createrole.UseCase
	updateRole         updaterole.UseCase
	deleteRole         deleterole.UseCase
	getRoles           getroles.UseCase
	setRolePermissions setrolepermissions.UseCase
	getRolePermissions getrolepermissions.UseCase

	// Access management
	setUserRoles       setuserroles.UseCase
	getUserRoles       getuserroles.UseCase
	setUserPermissions setuserpermissions.UseCase
	getUserPermissions getuserpermissions.UseCase

	// Session management
	getMySessions        getmysessions.UseCase
	deleteMySession      deletemysession.UseCase
	cleanExpiredSessions cleanexpiredsessions.UseCase
	getUserSessions      getusersessions.UseCase
	deleteSession        deletesession.UseCase
	deleteUserSessions   deleteusersessions.UseCase
}

//nolint:dupl // constructor args intentionally mirror container fields for DI clarity
func NewContainer(
	adminLogin adminlogin.UseCase,
	register register.UseCase,
	login login.UseCase,
	getMe getme.UseCase,
	googleOAuthLogin googleoauthlogin.UseCase,
	gitHubOAuthLogin githuboauthlogin.UseCase,
	refreshToken refreshtoken.UseCase,
	logout logout.UseCase,
	changeMyPassword changemypassword.UseCase,
	createSuperadmin createsuperadmin.UseCase,
	getUsers getusers.UseCase,
	getAuthStats getauthstats.UseCase,
	createUser createuser.UseCase,
	updateUser updateuser.UseCase,
	disableUser disableuser.UseCase,
	enableUser enableuser.UseCase,
	createRole createrole.UseCase,
	updateRole updaterole.UseCase,
	deleteRole deleterole.UseCase,
	getRoles getroles.UseCase,
	setRolePermissions setrolepermissions.UseCase,
	getRolePermissions getrolepermissions.UseCase,
	setUserRoles setuserroles.UseCase,
	getUserRoles getuserroles.UseCase,
	setUserPermissions setuserpermissions.UseCase,
	getUserPermissions getuserpermissions.UseCase,
	getMySessions getmysessions.UseCase,
	deleteMySession deletemysession.UseCase,
	cleanExpiredSessions cleanexpiredsessions.UseCase,
	getUserSessions getusersessions.UseCase,
	deleteSession deletesession.UseCase,
	deleteUserSessions deleteusersessions.UseCase,
) *Container {
	return &Container{
		adminLogin:           adminLogin,
		register:             register,
		login:                login,
		getMe:                getMe,
		googleOAuthLogin:     googleOAuthLogin,
		gitHubOAuthLogin:     gitHubOAuthLogin,
		refreshToken:         refreshToken,
		logout:               logout,
		changeMyPassword:     changeMyPassword,
		createSuperadmin:     createSuperadmin,
		getUsers:             getUsers,
		getAuthStats:         getAuthStats,
		createUser:           createUser,
		updateUser:           updateUser,
		disableUser:          disableUser,
		enableUser:           enableUser,
		createRole:           createRole,
		updateRole:           updateRole,
		deleteRole:           deleteRole,
		getRoles:             getRoles,
		setRolePermissions:   setRolePermissions,
		getRolePermissions:   getRolePermissions,
		setUserRoles:         setUserRoles,
		getUserRoles:         getUserRoles,
		setUserPermissions:   setUserPermissions,
		getUserPermissions:   getUserPermissions,
		getMySessions:        getMySessions,
		deleteMySession:      deleteMySession,
		cleanExpiredSessions: cleanExpiredSessions,
		getUserSessions:      getUserSessions,
		deleteSession:        deleteSession,
		deleteUserSessions:   deleteUserSessions,
	}
}

// --- Self-service ---

func (c *Container) AdminLogin() adminlogin.UseCase             { return c.adminLogin }
func (c *Container) Register() register.UseCase                 { return c.register }
func (c *Container) Login() login.UseCase                       { return c.login }
func (c *Container) GetMe() getme.UseCase                       { return c.getMe }
func (c *Container) GoogleOAuthLogin() googleoauthlogin.UseCase { return c.googleOAuthLogin }
func (c *Container) GitHubOAuthLogin() githuboauthlogin.UseCase { return c.gitHubOAuthLogin }
func (c *Container) RefreshToken() refreshtoken.UseCase         { return c.refreshToken }
func (c *Container) Logout() logout.UseCase                     { return c.logout }
func (c *Container) ChangeMyPassword() changemypassword.UseCase { return c.changeMyPassword }

// --- User management ---

func (c *Container) CreateSuperadmin() createsuperadmin.UseCase { return c.createSuperadmin }
func (c *Container) GetUsers() getusers.UseCase                 { return c.getUsers }
func (c *Container) GetAuthStats() getauthstats.UseCase         { return c.getAuthStats }
func (c *Container) CreateUser() createuser.UseCase             { return c.createUser }
func (c *Container) UpdateUser() updateuser.UseCase             { return c.updateUser }
func (c *Container) DisableUser() disableuser.UseCase           { return c.disableUser }
func (c *Container) EnableUser() enableuser.UseCase             { return c.enableUser }

// --- Role management ---

func (c *Container) CreateRole() createrole.UseCase                 { return c.createRole }
func (c *Container) UpdateRole() updaterole.UseCase                 { return c.updateRole }
func (c *Container) DeleteRole() deleterole.UseCase                 { return c.deleteRole }
func (c *Container) GetRoles() getroles.UseCase                     { return c.getRoles }
func (c *Container) SetRolePermissions() setrolepermissions.UseCase { return c.setRolePermissions }
func (c *Container) GetRolePermissions() getrolepermissions.UseCase { return c.getRolePermissions }

// --- Access management ---

func (c *Container) SetUserRoles() setuserroles.UseCase             { return c.setUserRoles }
func (c *Container) GetUserRoles() getuserroles.UseCase             { return c.getUserRoles }
func (c *Container) SetUserPermissions() setuserpermissions.UseCase { return c.setUserPermissions }
func (c *Container) GetUserPermissions() getuserpermissions.UseCase { return c.getUserPermissions }

// --- Session management ---

func (c *Container) GetMySessions() getmysessions.UseCase     { return c.getMySessions }
func (c *Container) DeleteMySession() deletemysession.UseCase { return c.deleteMySession }
func (c *Container) CleanExpiredSessions() cleanexpiredsessions.UseCase {
	return c.cleanExpiredSessions
}
func (c *Container) GetUserSessions() getusersessions.UseCase       { return c.getUserSessions }
func (c *Container) DeleteSession() deletesession.UseCase           { return c.deleteSession }
func (c *Container) DeleteUserSessions() deleteusersessions.UseCase { return c.deleteUserSessions }
