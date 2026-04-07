package adminlogin

import (
	"context"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/session"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/uow"
	"github.com/jaxongir1006/hire-ready-api/internal/modules/auth/domain/user"
	"github.com/jaxongir1006/hire-ready-api/internal/portal"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/audit"
	"github.com/jaxongir1006/hire-ready-api/internal/portal/auth"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/sorter"
	"github.com/rise-and-shine/pkg/token"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

var (
	errIncorrectCreds = errx.New(
		"username or password is incorrect",
		errx.WithType(errx.T_Validation),
		errx.WithCode(user.CodeIncorrectCreds),
	)
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required" mask:"true"`
}

type Response struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

// UseCase implements "admin-login" user action.
type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	portalContainer *portal.Container,
	accessTokenTTL,
	refreshTokenTTL time.Duration,
	maxActiveSessions int,
) UseCase {
	return &usecase{
		domainContainer:   domainContainer,
		portalContainer:   portalContainer,
		accessTokenTTL:    accessTokenTTL,
		refreshTokenTTL:   refreshTokenTTL,
		maxActiveSessions: maxActiveSessions,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container

	accessTokenTTL    time.Duration
	refreshTokenTTL   time.Duration
	maxActiveSessions int
}

func (uc *usecase) OperationID() string { return "admin-login" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Find user by username
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		Username: &in.Username,
	})
	if errx.IsCodeIn(err, user.CodeUserNotFound) {
		return nil, errx.Wrap(errIncorrectCreds, errx.WithDetails(errx.D{"cause": "username"}))
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	ctx = context.WithValue(ctx, meta.ActorType, auth.ActorTypeUser)
	ctx = context.WithValue(ctx, meta.ActorID, u.ID)

	// Check if user is active
	if !u.IsActive {
		return nil, errx.Wrap(errIncorrectCreds, errx.WithDetails(errx.D{"cause": "is_active"}))
	}

	// Verify password hash
	if u.PasswordHash == nil {
		return nil, errx.Wrap(errIncorrectCreds, errx.WithDetails(errx.D{"cause": "password"}))
	}
	ok := hasher.Compare(in.Password, *u.PasswordHash)
	if !ok {
		return nil, errx.Wrap(errIncorrectCreds, errx.WithDetails(errx.D{"cause": "password"}))
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Enforce max active sessions limit (delete least recently used sessions if exceeded)
	err = uc.deleteExceededSessions(ctx, u, uow)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Create session record with tokens and meta info (IP, user_agent)
	s, err := uow.Session().Create(ctx, &session.Session{
		UserID:                u.ID,
		AccessToken:           token.NewOpaqueToken(),
		AccessTokenExpiresAt:  time.Now().Add(uc.accessTokenTTL),
		RefreshToken:          token.NewOpaqueToken(),
		RefreshTokenExpiresAt: time.Now().Add(uc.refreshTokenTTL),
		IPAddress:             meta.Find(ctx, meta.IPAddress),
		UserAgent:             meta.Find(ctx, meta.UserAgent),
		LastUsedAt:            time.Now(),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Update user's last_active_at timestamp
	u.LastActiveAt = lo.ToPtr(time.Now())
	_, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Record audit log
	err = uc.portalContainer.Audit().Log(uow.Lend(), audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: in,
	})
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

func (uc *usecase) deleteExceededSessions(ctx context.Context, u *user.User, uow uow.UnitOfWork) error {
	activeSessions, err := uow.Session().List(ctx, session.Filter{
		UserID:            &u.ID,
		OrderByLastUsedAt: lo.ToPtr(sorter.Asc),
	})
	if err != nil {
		return errx.Wrap(err)
	}
	sessionsToDelete := len(activeSessions) - uc.maxActiveSessions + 1

	if sessionsToDelete <= 0 {
		return nil
	}

	err = uow.Session().BulkDelete(ctx, activeSessions[:sessionsToDelete])
	return errx.Wrap(err)
}
