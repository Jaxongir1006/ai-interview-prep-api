package login

import (
	"context"
	"strings"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/sessionmanager"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/ucdef"
)

var errIncorrectCreds = errx.New(
	"email or password is incorrect",
	errx.WithType(errx.T_Validation),
	errx.WithCode(user.CodeIncorrectCreds),
)

type Request struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"  mask:"true"`
}

type Response struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
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

func (uc *usecase) OperationID() string { return "login" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	email := normalizeEmail(in.Email)

	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		Email: &email,
	})
	if errx.IsCodeIn(err, user.CodeUserNotFound) {
		return nil, errx.Wrap(errIncorrectCreds)
	}
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if !u.IsActive {
		return nil, errx.Wrap(errIncorrectCreds)
	}
	if u.PasswordHash == nil || !hasher.Compare(in.Password, *u.PasswordHash) {
		return nil, errx.Wrap(errIncorrectCreds)
	}

	ctx = context.WithValue(ctx, meta.ActorType, auth.ActorTypeUser)
	ctx = context.WithValue(ctx, meta.ActorID, u.ID)

	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	s, err := uc.sessionManager.CreateAuthenticatedSession(ctx, uow, u)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	err = uc.portalContainer.Audit().Log(uow.Lend(), audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: in,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
