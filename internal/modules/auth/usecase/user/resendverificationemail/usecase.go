package resendverificationemail

import (
	"context"
	"strings"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/emailverification"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type Response struct{}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	emailVerificationService *emailverification.Service,
) UseCase {
	return &usecase{
		domainContainer:          domainContainer,
		emailVerificationService: emailVerificationService,
	}
}

type usecase struct {
	domainContainer          *domain.Container
	emailVerificationService *emailverification.Service
}

func (uc *usecase) OperationID() string { return "resend-verification-email" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Normalize email
	email := normalizeEmail(in.Email)

	// Find user by email
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		Email: &email,
	})
	if errx.IsCodeIn(err, user.CodeUserNotFound) {
		return &Response{}, nil
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Check that user has password-based credentials
	if u.PasswordHash == nil {
		return &Response{}, nil
	}

	// Check that user is active
	if !u.IsActive {
		return &Response{}, nil
	}

	// Check that user is not already verified
	if u.IsVerified {
		return &Response{}, nil
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Expire previous unused email verification tokens for this user and email
	// Create fresh one-time email verification token
	verificationToken, err := uc.emailVerificationService.CreateToken(ctx, uow, u, email)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Send verification email with frontend verification URL and raw token
	err = uc.emailVerificationService.SendVerificationEmail(ctx, email, verificationToken)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return empty response
	return &Response{}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
