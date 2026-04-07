package statuschangelog

import (
	"time"

	"github.com/rise-and-shine/pkg/repogen"
)

type Filter struct {
	ID          *int64
	ActionLogID *int64
	EntityType  *string
	EntityID    *string
	TraceID     *string

	CreatedFrom *time.Time
	CreatedTo   *time.Time

	Cursor *int64
	Limit  *int
}

type Repo interface {
	repogen.Repo[StatusChangeLog, Filter]
}
