package profile

import (
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID              *int64
	UserID          *string
	TargetRole      *string
	ExperienceLevel *string
	IDs             []int64

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[CandidateProfile, Filter]
}
