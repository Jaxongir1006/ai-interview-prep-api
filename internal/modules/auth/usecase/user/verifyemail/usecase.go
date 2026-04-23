package verifyemail

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/emailverificationtoken"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/emailverification"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/sessionmanager"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	Token string `json:"token" validate:"required"`
}

type Response struct {
	UserID                string `json:"user_id"`
	Email                 string `json:"email"`
	IsVerified            bool   `json:"is_verified"`
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
	OnboardingRequired    bool   `json:"onboarding_required"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	portalContainer *portal.Container,
	sessionManager *sessionmanager.Service,
) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		portalContainer: portalContainer,
		sessionManager:  sessionManager,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
	sessionManager  *sessionmanager.Service
}

func (uc *usecase) OperationID() string { return "verify-email" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Hash the raw token
	tokenHash := emailverification.HashToken(in.Token)

	// Consume Redis-backed email verification token by token hash
	evt, err := uc.domainContainer.EmailVerificationTokenRepo().Consume(ctx, tokenHash)
	if errx.IsCodeIn(err, emailverificationtoken.CodeEmailVerificationTokenNotFound) {
		return nil, invalidTokenErr()
	}
	if err != nil {
		return nil, errx.Wrap(err)
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

	// Mark user as verified
	u.IsVerified = true
	u, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	ctx = context.WithValue(ctx, meta.ActorType, auth.ActorTypeUser)
	ctx = context.WithValue(ctx, meta.ActorID, u.ID)

	// Create authenticated session
	s, err := uc.sessionManager.CreateAuthenticatedSession(ctx, uow, u)
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

	onboardingRequired, err := uc.onboardingRequired(ctx, u.ID)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return verified user identity data and authenticated session
	return &Response{
		UserID:                u.ID,
		Email:                 evt.Email,
		IsVerified:            u.IsVerified,
		AccessToken:           s.AccessToken,
		AccessTokenExpiresAt:  s.AccessTokenExpiresAt.Format(time.RFC3339),
		RefreshToken:          s.RefreshToken,
		RefreshTokenExpiresAt: s.RefreshTokenExpiresAt.Format(time.RFC3339),
		OnboardingRequired:    onboardingRequired,
	}, nil
}

func (uc *usecase) onboardingRequired(ctx context.Context, userID string) (bool, error) {
	p, err := uc.portalContainer.Candidate().GetProfileByUserID(ctx, userID)
	if errx.IsCodeIn(err, candidateportal.CodeProfileNotFound) {
		return true, nil
	}
	if err != nil {
		return false, errx.Wrap(err)
	}

	return !p.OnboardingCompleted, nil
}

func invalidTokenErr() error {
	return errx.New(
		"email verification token is invalid",
		errx.WithType(errx.T_Validation),
		errx.WithCode(emailverificationtoken.CodeEmailVerificationTokenInvalid),
	)
}
