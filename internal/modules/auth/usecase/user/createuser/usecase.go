package createuser

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/audit"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/mask"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Request struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"        mask:"true"`
}

type UseCase = ucdef.UserAction[*Request, *user.User]

func New(
	domainContainer *domain.Container,
	portalContainer *portal.Container,
	hashingCost int,
) UseCase {
	return &usecase{
		domainContainer, portalContainer, hashingCost,
	}
}

type usecase struct {
	domainContainer *domain.Container
	portalContainer *portal.Container

	hashingCost int
}

func (uc *usecase) OperationID() string { return "create-user" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*user.User, error) {
	// Hash the password
	passwordHash, err := hasher.Hash(in.Password, hasher.WithCost(uc.hashingCost))
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Create user
	u, err := uow.User().Create(ctx, &user.User{
		ID:           uuid.NewString(),
		Username:     &in.Username,
		PasswordHash: &passwordHash,
		IsActive:     true,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, user.CodeUsernameConflict)
	}

	// Record audit log
	err = uc.portalContainer.Audit().Log(uow.Lend(), audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: mask.StructToOrdMap(in),
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
