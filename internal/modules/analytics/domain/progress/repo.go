package progress

import "github.com/rise-and-shine/pkg/repogen"

type SummaryFilter struct {
	ID     *int64
	UserID *string
	IDs    []int64
}

type TopicStatFilter struct {
	ID       *int64
	UserID   *string
	TopicKey *string
	IDs      []int64
}

type SummaryRepo interface {
	repogen.Repo[Summary, SummaryFilter]
}

type TopicStatRepo interface {
	repogen.Repo[TopicStat, TopicStatFilter]
}
