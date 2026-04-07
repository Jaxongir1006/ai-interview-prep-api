package alerterror

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

type Error struct {
	pg.BaseModel

	ID string `json:"id" bun:"id,pk"`

	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details"    bun:",type:jsonb"`
	Service   string            `json:"service"`
	Operation string            `json:"operation"`
	CreatedAt time.Time         `json:"created_at"`
	Alerted   bool              `json:"alerted"`
}
