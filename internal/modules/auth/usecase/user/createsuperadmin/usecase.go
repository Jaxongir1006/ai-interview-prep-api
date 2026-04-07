package createsuperadmin

import (
	"context"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Input struct {
	Username string
	Password string
}

type UseCase = ucdef.ManualCommand[*Input]

func New(
	domainContainer *domain.Container,
	hashingCost int,
) UseCase {
	return &usecase{
		domainContainer,
		hashingCost,
	}
}

type usecase struct {
	domainContainer *domain.Container

	hashingCost int
}

func (uc *usecase) OperationID() string { return "create-superadmin" }

func (uc *usecase) Execute(ctx context.Context, input *Input) error {
	// Hash the password
	passwordHash, err := hasher.Hash(input.Password, hasher.WithCost(uc.hashingCost))
	if err != nil {
		return errx.Wrap(err)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Create user
	u, err := uow.User().Create(ctx, &user.User{
		ID:           uuid.NewString(),
		Username:     &input.Username,
		PasswordHash: &passwordHash,
		IsActive:     true,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	// Assign all superadmin permissions explicitly
	for _, perm := range auth.SuperadminPermissions() {
		_, err = uow.UserPermission().Create(ctx, &rbac.UserPermission{
			UserID:     u.ID,
			Permission: perm,
		})
		if err != nil {
			return errx.Wrap(err)
		}
	}

	// Apply UOW
	err = uow.ApplyChanges()
	return errx.Wrap(err)
}
