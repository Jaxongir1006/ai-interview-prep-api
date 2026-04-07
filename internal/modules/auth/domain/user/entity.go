package user

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeUserNotFound        = "USER_NOT_FOUND"
	CodeUsernameConflict    = "USERNAME_CONFLICT"
	CodeIncorrectCreds      = "INCORRECT_CREDENTIALS"
	CodeUserAlreadyActive   = "USER_ALREADY_ACTIVE"
	CodeUserAlreadyDisabled = "USER_ALREADY_DISABLED"
)

type User struct {
	pg.BaseModel

	ID string `json:"id" bun:"id,pk"`

	Username     *string `json:"username"`
	PasswordHash *string `json:"-"`

	IsActive     bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at"`
}
