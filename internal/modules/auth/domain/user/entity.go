package user

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeUserNotFound        = "USER_NOT_FOUND"
	CodeUsernameConflict    = "USERNAME_CONFLICT"
	CodeEmailConflict       = "EMAIL_CONFLICT"
	CodePhoneNumberConflict = "PHONE_NUMBER_CONFLICT"
	CodeIncorrectCreds      = "INCORRECT_CREDENTIALS"
	CodeUserAlreadyActive   = "USER_ALREADY_ACTIVE"
	CodeUserAlreadyDisabled = "USER_ALREADY_DISABLED"
)

type User struct {
	pg.BaseModel

	ID string `json:"id" bun:"id,pk"`

	Username     *string `json:"username"`
	Email        *string `json:"-"`
	PhoneNumber  *string `json:"-"`
	PasswordHash *string `json:"-"`
	IsVerified   bool    `json:"-"`

	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"-"`
	LastActiveAt *time.Time `json:"last_active_at"`
}
