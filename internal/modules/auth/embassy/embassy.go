package embassy

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
)

type embassy struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) auth.Portal {
	return &embassy{domainContainer: domainContainer}
}

func (e *embassy) Authenticate(ctx context.Context, accessToken string) (*auth.UserContext, error) {
	// Find session by access token
	s, err := e.domainContainer.SessionRepo().Get(ctx, session.Filter{
		AccessToken: &accessToken,
	})
	if errx.IsCodeIn(err, session.CodeSessionNotFound) {
		return nil, errx.New(
			"session not found",
			errx.WithType(errx.T_Authentication),
			errx.WithCode(auth.CodeUnauthorized),
		)
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Check access token expiration
	if time.Now().After(s.AccessTokenExpiresAt) {
		return nil, errx.New(
			"access token has expired",
			errx.WithType(errx.T_Authentication),
			errx.WithCode(auth.CodeSessionExpired),
		)
	}

	// Find user by session user ID
	u, err := e.domainContainer.UserRepo().Get(ctx, user.Filter{
		ID: &s.UserID,
	})
	if err != nil {
		return nil, errx.New(
			"user not found for session",
			errx.WithType(errx.T_Authentication),
			errx.WithCode(auth.CodeUnauthorized),
		)
	}

	// Check if user is active
	if !u.IsActive {
		return nil, errx.New(
			"user account is inactive",
			errx.WithType(errx.T_Authentication),
			errx.WithCode(auth.CodeUserInactive),
		)
	}

	// Collect user permissions
	permissions, err := e.domainContainer.UserPermissionRepo().CollectUserPermissions(ctx, u.ID)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &auth.UserContext{
		UserID:      u.ID,
		Username:    u.Username,
		IsVerified:  u.IsVerified,
		SessionID:   s.ID,
		Permissions: permissions,
	}, nil
}
