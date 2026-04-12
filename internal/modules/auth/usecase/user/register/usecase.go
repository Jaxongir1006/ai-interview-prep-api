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

type Request struct {
	Email    string `json:"email"     validate:"required,email,max=255"`
	FullName string `json:"full_name" validate:"required,min=1,max=255"`
	Password string `json:"password"  validate:"required,min=8,max=72"  mask:"true"`
}

type Response struct {
	User                 UserInfo    `json:"user"`
	Profile              ProfileInfo `json:"profile"`
	VerificationRequired bool        `json:"verification_required"`
}

type UserInfo struct {
	ID           string  `json:"id"`
	Email        *string `json:"email"`
	PhoneNumber  *string `json:"phone_number"`
	IsVerified   bool    `json:"is_verified"`
	IsActive     bool    `json:"is_active"`
	LastLoginAt  any     `json:"last_login_at"`
	LastActiveAt any     `json:"last_active_at"`
	CreatedAt    any     `json:"created_at"`
	UpdatedAt    any     `json:"updated_at"`
}

type ProfileInfo struct {
	ID                   int64    `json:"id"`
	UserID               string   `json:"user_id"`
	FullName             *string  `json:"full_name"`
	Bio                  *string  `json:"bio"`
	Location             *string  `json:"location"`
	TargetRole           *string  `json:"target_role"`
	ExperienceLevel      *string  `json:"experience_level"`
	InterviewGoalPerWeek int      `json:"interview_goal_per_week"`
	PreferredTopics      []string `json:"preferred_topics"`
	CreatedAt            any      `json:"created_at"`
	UpdatedAt            any      `json:"updated_at"`
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
	p, err := uc.portalContainer.Candidate().
		CreateInitialProfile(uow.Lend(), &candidateportal.CreateInitialProfileRequest{
			UserID:   u.ID,
			FullName: &in.FullName,
		})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Create one-time email verification token for the user's email
	verificationToken, err := uc.emailVerificationService.CreateToken(ctx, uow, u, email)
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

	// Apply UOW
	err = uow.ApplyChanges()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Send verification email with frontend verification URL and raw token
	err = uc.emailVerificationService.SendVerificationEmail(ctx, email, verificationToken)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Return created user and profile data with verification_required true
	return &Response{
		User: UserInfo{
			ID:           u.ID,
			Email:        u.Email,
			PhoneNumber:  u.PhoneNumber,
			IsVerified:   u.IsVerified,
			IsActive:     u.IsActive,
			LastLoginAt:  u.LastLoginAt,
			LastActiveAt: u.LastActiveAt,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
		},
		Profile: ProfileInfo{
			ID:                   p.ID,
			UserID:               p.UserID,
			FullName:             p.FullName,
			Bio:                  p.Bio,
			Location:             p.Location,
			TargetRole:           p.TargetRole,
			ExperienceLevel:      p.ExperienceLevel,
			InterviewGoalPerWeek: p.InterviewGoalPerWeek,
			PreferredTopics:      []string{},
			CreatedAt:            p.CreatedAt,
			UpdatedAt:            p.UpdatedAt,
		},
		VerificationRequired: true,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
