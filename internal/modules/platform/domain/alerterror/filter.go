package alerterror

import (
	"time"

	"github.com/rise-and-shine/pkg/sorter"
)

type Filter struct {
	ID        *string
	Code      *string
	Service   *string
	Operation *string
	Alerted   *bool

	CreatedFrom *time.Time
	CreatedTo   *time.Time

	Search string

	Limit    *int
	Offset   *int
	SortOpts sorter.SortOpts
}
