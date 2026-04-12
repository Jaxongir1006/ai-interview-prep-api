package emailverification

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	domainmail "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/mail"
	authuow "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/uow"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/token"
)

type Service struct {
	tokenTTL          time.Duration
	frontendVerifyURL string
	mailSender        domainmail.Sender
}

func New(
	tokenTTL time.Duration,
	frontendVerifyURL string,
	mailSender domainmail.Sender,
) *Service {
	return &Service{
		tokenTTL:          tokenTTL,
		frontendVerifyURL: frontendVerifyURL,
		mailSender:        mailSender,
	}
}

type CreatedToken struct {
	RawToken string
	URL      string
}

func (s *Service) CreateToken(
	ctx context.Context,
	uow authuow.UnitOfWork,
	u *user.User,
	email string,
) (*CreatedToken, error) {
	now := time.Now()

	err := uow.EmailVerificationToken().ExpireUnused(ctx, u.ID, email, now)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	rawToken := token.NewOpaqueToken()
	tokenHash := HashToken(rawToken)

	_, err = uow.EmailVerificationToken().Create(ctx, &emailverificationtoken.EmailVerificationToken{
		UserID:    u.ID,
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(s.tokenTTL),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &CreatedToken{
		RawToken: rawToken,
		URL:      s.VerificationURL(rawToken),
	}, nil
}

func (s *Service) SendVerificationEmail(
	ctx context.Context,
	email string,
	createdToken *CreatedToken,
) error {
	err := s.mailSender.SendVerificationEmail(ctx, domainmail.VerificationEmail{
		To:              email,
		VerificationURL: createdToken.URL,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func (s *Service) VerificationURL(rawToken string) string {
	u, err := url.Parse(s.frontendVerifyURL)
	if err != nil {
		return s.frontendVerifyURL
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
