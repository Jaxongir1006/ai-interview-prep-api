package sessionmanager

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	authuow "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/sorter"
	"github.com/rise-and-shine/pkg/token"
	"github.com/samber/lo"
)

type Service struct {
	accessTokenTTL    time.Duration
	refreshTokenTTL   time.Duration
	maxActiveSessions int
}

func New(
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	maxActiveSessions int,
) *Service {
	return &Service{
		accessTokenTTL:    accessTokenTTL,
		refreshTokenTTL:   refreshTokenTTL,
		maxActiveSessions: maxActiveSessions,
	}
}

func (s *Service) CreateAuthenticatedSession(
	ctx context.Context,
	uow authuow.UnitOfWork,
	u *user.User,
) (*session.Session, error) {
	err := s.deleteExceededSessions(ctx, uow, u.ID)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	now := time.Now()
	sess, err := uow.Session().Create(ctx, &session.Session{
		UserID:                u.ID,
		AccessToken:           token.NewOpaqueToken(),
		AccessTokenExpiresAt:  now.Add(s.accessTokenTTL),
		RefreshToken:          token.NewOpaqueToken(),
		RefreshTokenExpiresAt: now.Add(s.refreshTokenTTL),
		IPAddress:             meta.Find(ctx, meta.IPAddress),
		UserAgent:             meta.Find(ctx, meta.UserAgent),
		LastUsedAt:            now,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	u.LastLoginAt = lo.ToPtr(now)
	u.LastActiveAt = lo.ToPtr(now)
	_, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return sess, nil
}

func (s *Service) deleteExceededSessions(
	ctx context.Context,
	uow authuow.UnitOfWork,
	userID string,
) error {
	activeSessions, err := uow.Session().List(ctx, session.Filter{
		UserID:            &userID,
		OrderByLastUsedAt: lo.ToPtr(sorter.Asc),
	})
	if err != nil {
		return errx.Wrap(err)
	}

	sessionsToDelete := len(activeSessions) - s.maxActiveSessions + 1
	if sessionsToDelete <= 0 {
		return nil
	}

	err = uow.Session().BulkDelete(ctx, activeSessions[:sessionsToDelete])
	return errx.Wrap(err)
}
