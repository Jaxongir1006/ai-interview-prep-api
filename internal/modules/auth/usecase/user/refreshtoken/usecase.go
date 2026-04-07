package refreshtoken

import (
	"context"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/token"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Response struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

// UseCase implements "refresh-token" user action.
type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	accessTokenTTL,
	refreshTokenTTL time.Duration,
) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type usecase struct {
	domainContainer *domain.Container
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (uc *usecase) OperationID() string { return "refresh-token" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find session by refresh token
	s, err := uc.domainContainer.SessionRepo().Get(ctx, session.Filter{
		RefreshToken: &in.RefreshToken,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_Validation, session.CodeSessionNotFound)
	}

	// Check if refresh token is not expired
	if time.Now().After(s.RefreshTokenExpiresAt) {
		return nil, errx.New(
			"refresh token has expired",
			errx.WithType(errx.T_Validation),
			errx.WithCode(session.CodeSessionNotFound),
		)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Generate new access and refresh tokens with updated expiry
	s.AccessToken = token.NewOpaqueToken()
	s.AccessTokenExpiresAt = time.Now().Add(uc.accessTokenTTL)
	s.RefreshToken = token.NewOpaqueToken()
	s.RefreshTokenExpiresAt = time.Now().Add(uc.refreshTokenTTL)

	// Update session record with new tokens and meta info (IP, user_agent)
	s.IPAddress = meta.Find(ctx, meta.IPAddress)
	s.UserAgent = meta.Find(ctx, meta.UserAgent)
	s.LastUsedAt = time.Now()
	_, err = uow.Session().Update(ctx, s)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Update user's last_active_at timestamp
	u, err := uow.User().Get(ctx, userFilter(s.UserID))
	if err != nil {
		return nil, errx.Wrap(err)
	}
	u.LastActiveAt = lo.ToPtr(time.Now())
	_, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{
		AccessToken:           s.AccessToken,
		AccessTokenExpiresAt:  s.AccessTokenExpiresAt.Format(time.RFC3339),
		RefreshToken:          s.RefreshToken,
		RefreshTokenExpiresAt: s.RefreshTokenExpiresAt.Format(time.RFC3339),
	}, nil
}

func userFilter(id string) user.Filter {
	return user.Filter{ID: &id}
}
