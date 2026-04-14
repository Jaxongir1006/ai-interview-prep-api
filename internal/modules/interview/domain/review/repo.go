package review

import (
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID                 *int64
	SessionQuestionID  *int64
	AnswerID           *int64
	ReviewerType       *string
	IDs                []int64
	SessionQuestionIDs []int64
	AnswerIDs          []int64

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[Review, Filter]
}
