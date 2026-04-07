package actionlog

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

const (
	CodeActionLogNotFound = "ACTION_LOG_NOT_FOUND"
)

type ActionLog struct {
	ID             int64     `json:"id"              bun:"id,pk,autoincrement"`
	UserID         *string   `json:"user_id"`
	Module         string    `json:"module"`
	OperationID    string    `json:"operation_id"`
	RequestPayload any       `json:"request_payload" bun:",type:jsonb"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	Tags           []string  `json:"tags"            bun:",array,default:'{}'"`
	GroupKey       *string   `json:"group_key"`
	TraceID        string    `json:"trace_id"`
	CreatedAt      time.Time `json:"created_at"`
}

var _ bun.BeforeScanRowHook = (*ActionLog)(nil)

// BeforeScanRow resets the RequestPayload to nil before scanning so that
// bun's scanJSONIntoInterface takes the addressable nil-interface path
// instead of panicking on reflect.Value.Addr of a non-addressable Elem.
func (a *ActionLog) BeforeScanRow(_ context.Context) error {
	a.RequestPayload = nil
	return nil
}
