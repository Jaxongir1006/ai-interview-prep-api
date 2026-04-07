package user

import (
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID       *string
	Username *string
	IsActive *bool

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[User, Filter]
}
