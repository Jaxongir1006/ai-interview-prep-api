package changemypassword

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/ucdef"
)

var errIncorrectCreds = errx.New(
	"current password is incorrect",
	errx.WithType(errx.T_Validation),
	errx.WithCode(user.CodeIncorrectCreds),
)

type Request struct {
	CurrentPassword string `json:"current_password" validate:"required"       mask:"true"`
	NewPassword     string `json:"new_password"     validate:"required,min=8" mask:"true"`
}

type Response struct{}

// UseCase implements "change-my-password" user action.
type UseCase = ucdef.UserAction[*Request, *Response]

func New(domainContainer *domain.Container, portalContainer *portal.Container, hashingCost int) UseCase {
	return &usecase{
		domainContainer: domainContainer,
		portalContainer: portalContainer,
		hashingCost:     hashingCost,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container

	hashingCost int
}

func (uc *usecase) OperationID() string { return "change-my-password" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Get user context from authenticated session
	userCtx := auth.MustUserContext(ctx)

	// Find user by ID
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &userCtx.UserID})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Verify current password matches stored hash
	if u.PasswordHash == nil || !hasher.Compare(in.CurrentPassword, *u.PasswordHash) {
		return nil, errx.Wrap(errIncorrectCreds)
	}

	// Hash new password
	passwordHash, err := hasher.Hash(in.NewPassword, hasher.WithCost(uc.hashingCost))
	if err != nil {
		return nil, errx.Wrap(err)
	}
	u.PasswordHash = &passwordHash

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

	// Record audit log
	err = uc.portalContainer.Audit().Log(uow.Lend(), audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: in,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Response{}, nil
}
