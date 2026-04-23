package confirmpasswordreset

import (
	"context"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/passwordresettoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/session"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/passwordreset"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/ucdef"
)

const successMessage = "Your password has been reset. Sign in with your new password."

type Request struct {
	Token    string `json:"token"    validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=72" mask:"true"`
}

type Response struct {
	Message string `json:"message"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, hashingCost int) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		hashingCost:     hashingCost,
	}
}

type usecase struct {
	domainContainer *domain.Container
	hashingCost     int
}

func (uc *usecase) OperationID() string { return "confirm-password-reset" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Hash the raw token
	tokenHash := passwordreset.HashToken(in.Token)

	// Consume Redis-backed password reset token by token hash
	resetToken, err := uc.domainContainer.PasswordResetTokenRepo().Consume(ctx, tokenHash)
	if errx.IsCodeIn(err, passwordresettoken.CodePasswordResetTokenNotFound) {
		return nil, invalidTokenErr()
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Find user by token user ID
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		ID: &resetToken.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Check that user is active
	if !u.IsActive {
		return nil, invalidTokenErr()
	}

	// Check that user has password-based credentials
	if u.PasswordHash == nil {
		return nil, invalidTokenErr()
	}

	// Check that token email matches the user's current email
	if u.Email == nil || *u.Email != resetToken.Email {
		return nil, invalidTokenErr()
	}

	// Hash the new password
	passwordHash, err := hasher.Hash(in.Password, hasher.WithCost(uc.hashingCost))
	if err != nil {
		return nil, errx.Wrap(err)
	}
	u.PasswordHash = &passwordHash
	u.IsVerified = true

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Update user's password_hash
	_, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Mark user as verified
	// Already included in the user update.

	// Delete all sessions for the user
	sessions, err := uow.Session().List(ctx, session.Filter{
		UserID: &u.ID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if len(sessions) > 0 {
		err = uow.Session().BulkDelete(ctx, sessions)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return success response
	return &Response{Message: successMessage}, nil
}

func invalidTokenErr() error {
	return errx.New(
		"password reset token is invalid",
		errx.WithType(errx.T_Validation),
		errx.WithCode(passwordresettoken.CodePasswordResetTokenInvalid),
	)
}
