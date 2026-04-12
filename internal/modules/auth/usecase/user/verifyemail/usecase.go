package verifyemail

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/emailverification"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	Token string `json:"token" validate:"required"`
}

type Response struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, portalContainer *portal.Container) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		portalContainer: portalContainer,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
}

func (uc *usecase) OperationID() string { return "verify-email" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Hash the raw token
	tokenHash := emailverification.HashToken(in.Token)

	// Find unused email verification token by token hash
	unused := true
	evt, err := uc.domainContainer.EmailVerificationTokenRepo().Get(ctx, emailverificationtoken.Filter{
		TokenHash: &tokenHash,
		Unused:    &unused,
	})
	if errx.IsCodeIn(err, emailverificationtoken.CodeEmailVerificationTokenNotFound) {
		return nil, invalidTokenErr()
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Check that token is not expired
	if time.Now().After(evt.ExpiresAt) {
		return nil, errx.New(
			"email verification token is expired",
			errx.WithType(errx.T_Validation),
			errx.WithCode(emailverificationtoken.CodeEmailVerificationTokenExpired),
		)
	}

	// Find user by token user ID
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		ID: &evt.UserID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Check that token email matches the user's current email
	if u.Email == nil || *u.Email != evt.Email {
		return nil, errx.New(
			"email verification token does not match current user email",
			errx.WithType(errx.T_Validation),
			errx.WithCode(emailverificationtoken.CodeEmailVerificationEmailMismatch),
		)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Mark email verification token as used
	now := time.Now()
	evt.UsedAt = &now
	_, err = uow.EmailVerificationToken().Update(ctx, evt)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Mark user as verified
	u.IsVerified = true
	u, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Record audit log
	auditCtx := context.WithValue(uow.Lend(), meta.ActorType, auth.ActorTypeUser)
	auditCtx = context.WithValue(auditCtx, meta.ActorID, u.ID)
	err = uc.portalContainer.Audit().Log(auditCtx, audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: map[string]string{
			"user_id": u.ID,
			"email":   evt.Email,
		},
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return verified user identity data
	return &Response{
		UserID:     u.ID,
		Email:      evt.Email,
		IsVerified: u.IsVerified,
	}, nil
}

func invalidTokenErr() error {
	return errx.New(
		"email verification token is invalid",
		errx.WithType(errx.T_Validation),
		errx.WithCode(emailverificationtoken.CodeEmailVerificationTokenInvalid),
	)
}
