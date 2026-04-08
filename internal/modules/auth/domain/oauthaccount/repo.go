package oauthaccount

import (
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID             *int64
	UserID         *string
	Provider       *string
	ProviderUserID *string
	ProviderEmail  *string

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[OAuthAccount, Filter]
}
