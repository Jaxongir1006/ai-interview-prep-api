package passwordreset

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"time"

	domainmail "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/mail"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/passwordresettoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/token"
)

type Service struct {
	tokenTTL         time.Duration
	frontendResetURL string
	tokenRepo        passwordresettoken.Repo
	mailSender       domainmail.Sender
}

func New(
	tokenTTL time.Duration,
	frontendResetURL string,
	tokenRepo passwordresettoken.Repo,
	mailSender domainmail.Sender,
) *Service {
	return &Service{
		tokenTTL:         tokenTTL,
		frontendResetURL: frontendResetURL,
		tokenRepo:        tokenRepo,
		mailSender:       mailSender,
	}
}

type CreatedToken struct {
	RawToken string
	URL      string
}

func (s *Service) CreateToken(
	ctx context.Context,
	u *user.User,
	email string,
) (*CreatedToken, error) {
	now := time.Now()

	err := s.tokenRepo.InvalidateUserEmail(ctx, u.ID, email)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	rawToken := token.NewOpaqueToken()
	tokenHash := HashToken(rawToken)

	err = s.tokenRepo.Create(ctx, tokenHash, &passwordresettoken.PasswordResetToken{
		UserID:    u.ID,
		Email:     email,
		ExpiresAt: now.Add(s.tokenTTL),
	}, s.tokenTTL)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &CreatedToken{
		RawToken: rawToken,
		URL:      s.ResetURL(rawToken),
	}, nil
}

func (s *Service) SendPasswordResetEmail(
	ctx context.Context,
	email string,
	createdToken *CreatedToken,
) error {
	err := s.mailSender.SendPasswordResetEmail(ctx, domainmail.PasswordResetEmail{
		To:       email,
		ResetURL: createdToken.URL,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func (s *Service) ResetURL(rawToken string) string {
	u, err := url.Parse(s.frontendResetURL)
	if err != nil {
		return s.frontendResetURL
	}

	q := u.Query()
	q.Set("token", rawToken)
	u.RawQuery = q.Encode()

	return u.String()
}

func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
