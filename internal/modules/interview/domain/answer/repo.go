package answer

import (
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID                 *int64
	SessionQuestionID  *int64
	IDs                []int64
	SessionQuestionIDs []int64

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[Answer, Filter]
}
