package githuboauthlogin

import (
	"context"
	"time"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/oauthaccount"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	authoauth "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/infra/oauth"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/sessionmanager"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"

	"github.com/code19m/errx"
	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/ucdef"
	"github.com/samber/lo"
)

type Request struct {
	Code string `json:"code" validate:"required"`
}

type Response struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
	IsNewUser             bool   `json:"is_new_user"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	portalContainer *portal.Container,
	sessionManager *sessionmanager.Service,
	provider authoauth.GitHubProvider,
) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		portalContainer: portalContainer,
		sessionManager:  sessionManager,
		provider:        provider,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container
	sessionManager  *sessionmanager.Service
	provider        authoauth.GitHubProvider
}

func (uc *usecase) OperationID() string { return "github-oauth-login" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	identity, err := uc.provider.AuthenticateCode(ctx, in.Code)
	if err != nil {
		return nil, errx.Wrap(err, errx.WithType(errx.T_Validation))
	}

	oauthAcc, userEntity, err := uc.resolveAccount(ctx, identity)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if userEntity != nil && !userEntity.IsActive {
		return nil, errx.New("user account is inactive",
			errx.WithType(errx.T_Validation),
			errx.WithCode(auth.CodeUserInactive),
		)
	}

	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	isNewUser := false
	if userEntity == nil {
		userEntity, err = uow.User().Create(ctx, &user.User{
			ID:         uuid.NewString(),
			Email:      lo.ToPtr(identity.Email),
			IsVerified: identity.EmailVerified,
			IsActive:   true,
		})
		if err != nil {
			return nil, errx.Wrap(err)
		}
		isNewUser = true
	}

	now := time.Now()
	if oauthAcc == nil {
		_, err = uow.OAuthAccount().Create(ctx, &oauthaccount.OAuthAccount{
			UserID:         userEntity.ID,
			Provider:       oauthaccount.ProviderGitHub,
			ProviderUserID: identity.ProviderUserID,
			ProviderEmail:  lo.ToPtr(identity.Email),
			LastLoginAt:    lo.ToPtr(now),
		})
		if err != nil {
			return nil, errx.Wrap(err)
		}
	} else {
		oauthAcc.ProviderEmail = lo.ToPtr(identity.Email)
		oauthAcc.LastLoginAt = lo.ToPtr(now)
		_, err = uow.OAuthAccount().Update(ctx, oauthAcc)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	ctx = context.WithValue(ctx, meta.ActorType, auth.ActorTypeUser)
	ctx = context.WithValue(ctx, meta.ActorID, userEntity.ID)

	if isNewUser {
		_, err = uc.portalContainer.Candidate().
			CreateInitialProfile(uow.Lend(), &candidateportal.CreateInitialProfileRequest{
				UserID:   userEntity.ID,
				FullName: identity.FullName,
			})
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	s, err := uc.sessionManager.CreateAuthenticatedSession(ctx, uow, userEntity)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	auditCtx := context.WithValue(uow.Lend(), meta.ActorType, auth.ActorTypeUser)
	auditCtx = context.WithValue(auditCtx, meta.ActorID, userEntity.ID)

	err = uc.portalContainer.Audit().Log(auditCtx, audit.Action{
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
		IsNewUser:             isNewUser,
	}, nil
}

func (uc *usecase) resolveAccount(
	ctx context.Context,
	identity *authoauth.Identity,
) (*oauthaccount.OAuthAccount, *user.User, error) {
	oauthAcc, err := uc.domainContainer.OAuthAccountRepo().Get(ctx, oauthaccount.Filter{
		Provider:       lo.ToPtr(oauthaccount.ProviderGitHub),
		ProviderUserID: &identity.ProviderUserID,
	})
	if err == nil {
		u, userErr := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &oauthAcc.UserID})
		if userErr != nil {
			return nil, nil, errx.Wrap(userErr)
		}
		return oauthAcc, u, nil
	}
	if !errx.IsCodeIn(err, oauthaccount.CodeOAuthAccountNotFound) {
		return nil, nil, errx.Wrap(err)
	}

	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		Email: lo.ToPtr(identity.Email),
	})
	if errx.IsCodeIn(err, user.CodeUserNotFound) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, errx.Wrap(err)
	}

	return nil, u, nil
}
