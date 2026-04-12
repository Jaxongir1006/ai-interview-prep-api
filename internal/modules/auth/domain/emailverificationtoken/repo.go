package emailverificationtoken

import (
	"context"
	"time"

	"github.com/rise-and-shine/pkg/repogen"
)

type Filter struct {
	ID        *int64
	UserID    *string
	Email     *string
	TokenHash *string

	Unused *bool

	Limit *int
}

type Repo interface {
	repogen.Repo[EmailVerificationToken, Filter]

	ExpireUnused(ctx context.Context, userID, email string, now time.Time) error
}
