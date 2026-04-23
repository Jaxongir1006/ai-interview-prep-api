package register

import (
	"context"
	"strings"

	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/user"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/pblc/emailverification"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/audit"
	"github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/auth"
	candidateportal "github.com/Jaxongir1006/ai-interview-prep-api/internal/portal/candidate"

	"github.com/code19m/errx"
	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/mask"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/ucdef"
)

var errEmailConflict = errx.New(
	"email already exists",
	errx.WithType(errx.T_Conflict),
	errx.WithCode(user.CodeEmailConflict),
)

type Request struct {
	Email    string `json:"email"     validate:"required,email,max=255"`
	FullName string `json:"full_name" validate:"required,min=1,max=255"`
	Password string `json:"password"  validate:"required,min=8,max=72"  mask:"true"`
}

type Response struct {
	Email                string `json:"email"`
	VerificationRequired bool   `json:"verification_required"`
}

type UseCase = ucdef.UserAction[*Request, *Response]

func New(
	domainContainer *domain.Container,
	portalContainer *portal.Container,
	emailVerificationService *emailverification.Service,
	hashingCost int,
) UseCase {
	return &usecase{
		domainContainer:          domainContainer,
		portalContainer:          portalContainer,
		emailVerificationService: emailVerificationService,
		hashingCost:              hashingCost,
	}
}

type usecase struct {
	domainContainer          *domain.Container
	portalContainer          *portal.Container
	emailVerificationService *emailverification.Service
	hashingCost              int
}

func (uc *usecase) OperationID() string { return "register" }

func (uc *usecase) Execute(ctx context.Context, in *Request) (*Response, error) {
	// Normalize email
	email := normalizeEmail(in.Email)

	// Check whether a user already exists with the same email
	existingUser, err := uc.domainContainer.UserRepo().Get(ctx, user.Filter{
		Email: &email,
	})
	if err == nil {
		return uc.handleExistingUser(ctx, existingUser, email)
	}
	if !errx.IsCodeIn(err, user.CodeUserNotFound) {
		return nil, errx.Wrap(err)
	}

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

	// Create auth user with email-based credentials
	u, err := uow.User().Create(ctx, &user.User{
		ID:           uuid.NewString(),
		Email:        &email,
		PasswordHash: &passwordHash,
		IsActive:     true,
	})
	if err != nil {
		return nil, errx.WrapWithTypeOnCodes(err, errx.T_Conflict, user.CodeEmailConflict)
	}

	// Set user is_verified to false
	u.IsVerified = false

	// Create minimal candidate profile for the new user using the provided full name
	_, err = uc.portalContainer.Candidate().
		CreateInitialProfile(uow.Lend(), &candidateportal.CreateInitialProfileRequest{
			UserID:   u.ID,
			FullName: &in.FullName,
		})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Record audit log
	auditCtx := context.WithValue(uow.Lend(), meta.ActorType, auth.ActorTypeUser)
	auditCtx = context.WithValue(auditCtx, meta.ActorID, u.ID)

	err = uc.portalContainer.Audit().Log(auditCtx, audit.Action{
		Module: auth.ModuleName, OperationID: uc.OperationID(), Payload: mask.StructToOrdMap(in),
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Create fresh Redis-backed one-time email verification token for the user's email
	verificationToken, err := uc.emailVerificationService.CreateToken(ctx, u, email)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return uc.sendVerificationEmail(ctx, email, verificationToken)
}

func (uc *usecase) handleExistingUser(
	ctx context.Context,
	u *user.User,
	email string,
) (*Response, error) {
	if u.PasswordHash == nil || !u.IsActive || u.IsVerified {
		return nil, errEmailConflict
	}

	// Create a fresh Redis-backed email verification token
	verificationToken, err := uc.emailVerificationService.CreateToken(ctx, u, email)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return uc.sendVerificationEmail(ctx, email, verificationToken)
}

func (uc *usecase) sendVerificationEmail(
	ctx context.Context,
	email string,
	verificationToken *emailverification.CreatedToken,
) (*Response, error) {
	err := uc.emailVerificationService.SendVerificationEmail(ctx, email, verificationToken)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return normalized email with verification_required true
	return &Response{
		Email:                email,
		VerificationRequired: true,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
