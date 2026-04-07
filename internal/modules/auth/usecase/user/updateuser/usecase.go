package updateuser

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

type Request struct {
	ID       string  `json:"id"       validate:"required"`
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Password *string `json:"password" validate:"omitempty,min=8"        mask:"true"`
}

type UseCase = ucdef.UserAction[*Request, *user.User]

func New(domainContainer *domain.Container, portalContainer *portal.Container, hashingCost int) UseCase {
	return &usecase{domainContainer: domainContainer, portalContainer: portalContainer, hashingCost: hashingCost}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container

	hashingCost int
}

func (uc *usecase) OperationID() string { return "update-user" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*user.User, error) {
	// Find user by ID
	u, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{ID: &in.ID})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_NotFound, user.CodeUserNotFound)
	}

	// If password provided, hash the new password
	if in.Password != nil {
		passwordHash, hashErr := hasher.Hash(*in.Password, hasher.WithCost(uc.hashingCost))
		if hashErr != nil {
			return nil, errx.Wrap(hashErr)
		}
		u.PasswordHash = &passwordHash
	}

	// Update user fields
	if in.Username != nil {
		u.Username = in.Username
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Update user fields
	u, err = uow.User().Update(ctx, u)
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, user.CodeUsernameConflict)
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

	return u, nil
}
