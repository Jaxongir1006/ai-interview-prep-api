package statuschangelog

import "time"

const (
	CodeStatusChangeLogNotFound = "STATUS_CHANGE_LOG_NOT_FOUND"
)

type StatusChangeLog struct {
	ID          int64     `json:"id"            bun:"id,pk,autoincrement"`
	ActionLogID int64     `json:"action_log_id"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	Status      string    `json:"status"`
	TraceID     string    `json:"trace_id"`
	CreatedAt   time.Time `json:"created_at"`
}
