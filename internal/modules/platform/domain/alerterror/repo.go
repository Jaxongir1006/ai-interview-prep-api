package alerterror

import (
	"context"
	"time"
)

type ServiceCount struct {
	Service string `json:"service"`
	Count   int64  `json:"count"`
}

type OperationCount struct {
	Operation string `json:"operation"`
	Count     int64  `json:"count"`
}

type CodeCount struct {
	Code  string `json:"code"`
	Count int64  `json:"count"`
}

type Stats struct {
	TotalCount  int64            `json:"total_count"`
	ByService   []ServiceCount   `json:"by_service"`
	ByOperation []OperationCount `json:"by_operation"`
	ByCode      []CodeCount      `json:"by_code"`
}

type StatsFilter struct {
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

type Repo interface {
	Get(ctx context.Context, id string) (*Error, error)
	ListWithCount(ctx context.Context, f Filter) ([]Error, int64, error)
	GetStats(ctx context.Context, f StatsFilter) (*Stats, error)
	DeleteOlderThan(ctx context.Context, before time.Time) (int64, error)
}
