package question

import (
	"github.com/rise-and-shine/pkg/repogen"
	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID           *int64
	SessionID    *string
	TopicKey     *string
	Difficulty   *string
	QuestionType *string
	Source       *string
	IDs          []int64
	SessionIDs   []string
	TopicKeys    []string

	Limit  *int
	Offset *int

	SortOpts sorter.SortOpts
}

type Repo interface {
	repogen.Repo[Question, Filter]
}
