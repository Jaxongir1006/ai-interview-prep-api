package requestpasswordreset

import (
	"context"
	"strings"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/passwordreset"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const successMessage = "If that email can receive password reset instructions, a reset link is on the way."

type Request struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type Response struct {
	Message string `json:"message"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	passwordResetService *passwordreset.Service,
) UseCase {
	return &usecase{
		domainContainer:      domainContainer,
		passwordResetService: passwordResetService,
	}
}

type usecase struct {
	domainContainer      *domain.Container
	passwordResetService *passwordreset.Service
}

func (uc *usecase) OperationID() string { return "request-password-reset" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Normalize email
	email := normalizeEmail(in.Email)

	// Apply password reset request rate limit by normalized email and client IP
	// Rate limiting is not wired yet.

	// Find user by email
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		Email: &email,
	})
	if errx.IsCodeIn(err, user.CodeUserNotFound) {
		return genericSuccess(), nil
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Check that user has password-based credentials
	if u.PasswordHash == nil {
		return genericSuccess(), nil
	}

	// Check that user is active
	if !u.IsActive {
		return genericSuccess(), nil
	}

	// Invalidate previous Redis-backed password reset token for this user and email
	// Create fresh Redis-backed one-time password reset token
	resetToken, err := uc.passwordResetService.CreateToken(ctx, u, email)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Send password reset email with frontend reset URL and raw token
	err = uc.passwordResetService.SendPasswordResetEmail(ctx, email, resetToken)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return generic success response
	return genericSuccess(), nil
}

func genericSuccess() *Response {
	return &Response{Message: successMessage}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
