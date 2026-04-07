package getauthstats

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

type Request struct{}

type Response struct {
	TotalUsers     int `json:"total_users"`
	TotalRoles     int `json:"total_roles"`
	ActiveSessions int `json:"active_sessions"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container) UseCase {
	return &usecase{domainContainer: domainContainer}
}

type usecase struct {
	domainContainer *domain.Container
}

func (uc *usecase) OperationID() string { return "get-auth-stats" }

func (uc *usecase) Execute(ctx context.Context, _ *Request) (*Response, error) {
	// Count total users
	totalUsers, err := uc.domainContainer.UserRepo().Count(ctx, user.Filter{})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Count total roles
	totalRoles, err := uc.domainContainer.RoleRepo().Count(ctx, rbac.RoleFilter{})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Count active sessions (where refresh token has not expired)
	activeSessions, err := uc.domainContainer.SessionRepo().Count(ctx, session.Filter{
		IsActive: lo.ToPtr(true),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return aggregated stats
	return &Response{
		TotalUsers:     totalUsers,
		TotalRoles:     totalRoles,
		ActiveSessions: activeSessions,
	}, nil
}
