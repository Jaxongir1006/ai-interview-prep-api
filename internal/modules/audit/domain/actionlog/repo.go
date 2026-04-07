package actionlog

import (
	"time"

	"github.com/rise-and-shine/pkg/repogen"
)

type Filter struct {
	ID          *int64
	UserID      *string
	Module      *string
	OperationID *string
	TraceID     *string
	Tags        []string
	GroupKey    *string

	CreatedFrom *time.Time
	CreatedTo   *time.Time

	Cursor *int64
	Limit  *int
}

type Repo interface {
	repogen.Repo[ActionLog, Filter]
}
