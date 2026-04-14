package session

import (
	"time"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID              *string
	UserID          *string
	Status          *string
	TargetRole      *string
	ExperienceLevel *string
	Difficulty      *string
	StartedAtFrom   *time.Time
	StartedAtTo     *time.Time
	IDs             []string
	Statuses        []string

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[Session, Filter]
}
