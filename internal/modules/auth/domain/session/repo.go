package session

import (
	"context"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID           *int64
	UserID       *string
	AccessToken  *string
	RefreshToken *string

	IsActive *bool // true = refresh_token_expires_at >= now, false = refresh_token_expires_at < now

	Limit  *int
	Offset *int

	OrderByLastUsedAt *sorter.SortDirection
	SortOpts          sorter.SortOpts
}

type Repo interface {
	repogen.Repo[Session, Filter]

	// DeleteExpired deletes sessions where the refresh token has expired.
	// Returns the number of deleted sessions.
	DeleteExpired(ctx context.Context) (int64, error)
}
